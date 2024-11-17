package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"practica/internal/models"
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

func (c *Cliente) WriteJSON(v interface{}) error {
	return c.Conn.WriteJSON(v)
}

var Clientes = make(map[string]*Cliente)
var Broadcast = make(chan models.Notificaciones_All)
var Mutex = &sync.Mutex{}

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

		// Obtener el ID del usuario desde la query
		ID_usuario := c.Query("ID_usuario")
		id, err := strconv.ParseUint(ID_usuario, 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		Clientes[fmt.Sprint(id)] = NuevoCliente(conn, uint(id))

		for {
			mt, mensaje, err := conn.ReadMessage() // Leer el mensaje del cliente
			if err != nil {                        // Manejar errores
				fmt.Println(err)
				return
			}

			var wsNotificacion models.Notificaciones_All
			if err := json.Unmarshal(mensaje, &wsNotificacion); err != nil {
				fmt.Println("Error al deserializar mensaje:", err)
				continue
			}
			fmt.Printf("%s: %s\n", wsNotificacion.Titulo, wsNotificacion.Mensaje)

			err = conn.WriteMessage(mt, mensaje)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func HandleMessages() {
	for {
		msg := <-Broadcast
		if conn, ok := Clientes[fmt.Sprint(msg.Id)]; ok {
			err := conn.WriteJSON(msg)
			if err != nil {
				conn.Conn.Close()
				delete(Clientes, fmt.Sprint(msg.Id))
			}
		}
	}
}
