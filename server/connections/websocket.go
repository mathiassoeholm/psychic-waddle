package connections

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketMessage = []byte

type Websocket struct {
	socket *websocket.Conn
	// To the web
	outgoing chan WebsocketMessage
	// From the web
	incoming chan WebsocketMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWebsocket(resultChannel chan<- Websocket, writer http.ResponseWriter, request *http.Request) {
	log.Println("serveSocket start")
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// We give the clients 60 seconds to read Ping frames.
	pongWait := 60 * time.Second

	socket.SetReadDeadline(time.Now().Add(pongWait))
	socket.SetPongHandler(func(string) error { socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	websocket := Websocket{socket: socket, incoming: make(chan WebsocketMessage), outgoing: make(chan WebsocketMessage)}

	go websocket.readPump()
	go websocket.writePump()

	resultChannel <- websocket
	log.Println("serveSocket end")
}

func (ws *Websocket) readPump() {
	defer ws.socket.Close()

	for {
		_, message, err := ws.socket.ReadMessage()
		if err != nil {
			fmt.Println("readPump error:", err)
			break
		}

		ws.incoming <- message
	}
}

func (ws *Websocket) writePump() {
	defer ws.socket.Close()

	for {
		select {
		case message := <-ws.outgoing:
			err := ws.socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				fmt.Println("writePump error:", err)
				break
			}
		}
	}
}
