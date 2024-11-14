package websocket

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Conn *websocket.Conn
	Id   uint `json:"id"`
}

type ComentarioNotificacion struct {
	Mensaje string `json:"mensaje"`
}

var Clientes = make(map[string]*Cliente)
var Broadcast = make(chan MensajeNotificacion)
var Mutex = &sync.Mutex{}

type MensajeNotificacion struct {
	ID_remitente string `json:"ID_remitente"`
	Contenido    string `json:"Contenido"`
}

func NuevoCliente(Conn *websocket.Conn, ID uint) *Cliente {
	return &Cliente{
		Conn: Conn,
		Id:   ID,
	}
}

func Handle_WebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener la ID del usuario desde la URL
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()
		for {
			mt, mensaje, err := conn.ReadMessage() // Leer el mensaje del cliente
			if err != nil {                        // Manejar errores
				fmt.Println(err)
				return
			}
			fmt.Printf("Recibido: %s\n", mensaje)
			err = conn.WriteMessage(mt, mensaje)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func NotificarClientes(ID_remitente string, contenido string, ID_cliente string) {
	id, err := strconv.ParseUint(ID_cliente, 10, 64)
	if err != nil {
		log.Printf("error parsing ID_cliente: %v", err)
		return
	}
	mensaje := MensajeNotificacion{Contenido: contenido}

	Mutex.Lock()
	defer Mutex.Unlock()

	for _, cliente := range Clientes {
		if cliente.Id == uint(id) {
			err := cliente.Conn.WriteJSON(mensaje)
			if err != nil {
				log.Printf("error sending message: %v", err)
			}
		}
	}
}

func BroadcastMessage(c *gin.Context) {
	var json struct {
		ID_remitente string `json:"id_remitente"`
		Contenido    string `json:"contenido"`
		ID_cliente   string `json:"id_cliente"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	NotificarClientes(json.ID_remitente, json.Contenido, json.ID_cliente)
	c.JSON(http.StatusOK, gin.H{"status": "notificación enviada"})
}
