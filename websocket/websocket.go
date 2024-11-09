package websocket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)

type MensajeNotificacion struct {
	// Define the fields of MensajeNotificacion here
	Mensaje string `json:"mensaje"`
}

var Broadcast = make(chan MensajeNotificacion)
var mutex = &sync.Mutex{}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Mensaje recibido: \n", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func handleMessages() {
	for {
		msg := <-Broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

func init() {
	go handleMessages()
}
