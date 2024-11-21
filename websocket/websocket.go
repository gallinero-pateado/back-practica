package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practica/internal/database"
	"practica/internal/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Cliente struct {
	conn        *websocket.Conn
	temas       map[uint]bool // IDs de temas
	comentarios map[uint]bool // IDs de comentarios
}

func (c *Cliente) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

var clientes = make(map[string]*Cliente)
var Broadcast = make(chan models.Notificaciones_All)
var clientesMutex = &sync.Mutex{}

func Handle_WebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ID_usuario := c.Query("ID_usuario")
		if ID_usuario == "" {
			fmt.Println("Error: ID_usuario está vacío")
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Error al establecer WebSocket:", err)
			return
		}
		defer conn.Close()

		cliente, err := CargarSuscripcionesCliente(ID_usuario)
		if err != nil {
			fmt.Println("Error al cargar suscripciones del cliente:", err)
			return
		}
		cliente.conn = conn

		clientesMutex.Lock()
		clientes[ID_usuario] = cliente
		clientesMutex.Unlock()

		for {
			_, mensaje, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error al leer mensaje WebSocket:", err)
				clientesMutex.Lock()
				delete(clientes, ID_usuario)
				clientesMutex.Unlock()
				return
			}

			var notificacion models.Notificaciones_All
			if err := json.Unmarshal(mensaje, &notificacion); err != nil {
				fmt.Println("Error al deserializar mensaje:", err)
				continue
			}

			if notificacion.Id_comentario != 0 {
				clientesMutex.Lock()
				clientes[ID_usuario].comentarios[notificacion.Id_comentario] = true
				clientesMutex.Unlock()
			}

			fmt.Printf("Comentario ID %d: %s\n", notificacion.Id_comentario, notificacion.Mensaje)
		}
	}
}

func HandleMessages() {
	for {
		msg := <-Broadcast
		clientesMutex.Lock()
		for _, cliente := range clientes {
			err := cliente.conn.WriteJSON(msg)
			if err != nil {
				cliente.conn.Close()
				delete(clientes, fmt.Sprint(msg.Id))
			}
		}
		clientesMutex.Unlock()
	}
}

func EnviarMensajeATemaOComentario(mensaje models.Notificaciones_All) {
	comentarioPadreID, err := ObtenerComentarioPadreID(mensaje.Id_comentario)
	if err != nil {
		fmt.Println("Error al obtener comentario padre:", err)
		return
	}

	clientesMutex.Lock()
	defer clientesMutex.Unlock()

	for _, cliente := range clientes {
		if cliente.temas[comentarioPadreID] || cliente.comentarios[mensaje.Id_comentario] {
			err := cliente.conn.WriteJSON(mensaje)
			if err != nil {
				fmt.Println("Error al enviar mensaje:", err)
			}
		}
	}
}

func CargarSuscripcionesCliente(ID_usuario string) (*Cliente, error) {
	var temas []models.Tema
	var comentarios []models.Comentario

	// Obtener los temas en los que ha participado el usuario
	if err := database.DB.Where("usuario_id = ?", ID_usuario).Find(&temas).Error; err != nil {
		return nil, err
	}

	// Obtener los comentarios en los que ha participado el usuario
	if err := database.DB.Where("usuario_id = ?", ID_usuario).Find(&comentarios).Error; err != nil {
		return nil, err
	}

	cliente := &Cliente{
		temas:       make(map[uint]bool),
		comentarios: make(map[uint]bool),
	}

	for _, tema := range temas {
		cliente.temas[tema.ID] = true
	}

	for _, comentario := range comentarios {
		cliente.comentarios[comentario.ID] = true
	}

	return cliente, nil
}

func ObtenerComentarioPadreID(comentarioID uint) (uint, error) {
	var comentario models.Comentario
	if err := database.DB.First(&comentario, comentarioID).Error; err != nil {
		return 0, err
	}
	return comentario.Comentario_padre_id, nil
}
