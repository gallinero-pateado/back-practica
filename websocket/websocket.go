package websocket

import (
	"log"
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

type Cliente struct {
	conn *websocket.Conn
	id   string
}

type MensajeNotificacion struct {
	Mensaje string `json:"mensaje"`
}

var Clientes = make(map[string]*Cliente)
var Broadcast = make(chan MensajeNotificacion)
var Mutex = &sync.Mutex{}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	ClienteID := r.URL.Query().Get("id")
	if ClienteID == "" {
		log.Println("Client ID is required")
		return
	}

	cliente := &Cliente{conn: ws, id: ClienteID}
	Mutex.Lock()
	Clientes[ClienteID] = cliente
	Mutex.Unlock()

	for {
		var msg MensajeNotificacion
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(Clientes, ClienteID)
			break
		}
		Broadcast <- msg
	}
}

func SendNotification(clienteID string, mensaje string) {
	Mutex.Lock()
	cliente, ok := Clientes[clienteID]
	Mutex.Unlock()
	if ok {
		notificacion := MensajeNotificacion{Mensaje: mensaje}
		err := cliente.conn.WriteJSON(notificacion)
		if err != nil {
			log.Printf("error: %v", err)
			cliente.conn.Close()
			Mutex.Lock()
			delete(Clientes, clienteID)
			Mutex.Unlock()
		}
	}
}
