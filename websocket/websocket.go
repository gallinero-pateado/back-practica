package websocket

import (
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

func HandleWebSocket(c *gin.Context) {
	// Obtener la ID del usuario desde la URL
	id := c.Param("id")

	// Actualizar la conexión HTTP a una conexión WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar a WebSocket"})
		return
	}

	// Crear un nuevo cliente con la conexión WebSocket y la ID
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	cliente := &Cliente{Conn: conn, Id: uint(idUint)}

	// Agregar el cliente a la lista de clientes conectados
	Mutex.Lock()
	Clientes[id] = cliente
	Mutex.Unlock()

	// Manejar la conexión WebSocket
	for {
		var msg MensajeNotificacion
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			Mutex.Lock()
			delete(Clientes, id)
			Mutex.Unlock()
			break
		}
		Broadcast <- msg
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
