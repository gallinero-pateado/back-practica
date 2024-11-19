package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	temas       map[string]bool
	comentarios map[string]bool
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

		clientesMutex.Lock()
		clientes[ID_usuario] = &Cliente{conn: conn, temas: make(map[string]bool), comentarios: make(map[string]bool)}
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

			if notificacion.Titulo != "" {
				clientesMutex.Lock()
				clientes[ID_usuario].temas[notificacion.Titulo] = true
				clientesMutex.Unlock()
			}

			if notificacion.Mensaje != "" {
				clientesMutex.Lock()
				clientes[ID_usuario].comentarios[notificacion.Mensaje] = true
				clientesMutex.Unlock()
			}

			fmt.Printf("%s: %s\n", notificacion.Titulo, notificacion.Mensaje)
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

func EnviarMensajeATemaOComentario(tema string, comentario string, mensaje models.Notificaciones_All) {
	clientesMutex.Lock()
	defer clientesMutex.Unlock()

	for _, cliente := range clientes {
		if cliente.temas[tema] || cliente.comentarios[comentario] {
			err := cliente.conn.WriteJSON(mensaje)
			if err != nil {
				fmt.Println("Error al enviar mensaje:", err)
			}
		}
	}
}
