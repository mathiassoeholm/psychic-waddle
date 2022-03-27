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
	id     string
	// To the web
	outgoing chan WebsocketMessage
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

	websocket := Websocket{socket: socket, outgoing: make(chan WebsocketMessage)}

	resultChannel <- websocket
	log.Println("serveSocket end")
}

func (ws *Websocket) accept(playerId string, messageChannel chan<- IncomingWebsocketMessage) {
	ws.id = playerId

	go ws.readPump(messageChannel)
	go ws.writePump()

}

func (ws *Websocket) readPump(messageChannel chan<- IncomingWebsocketMessage) {
	defer ws.socket.Close()

	for {
		_, message, err := ws.socket.ReadMessage()
		fmt.Println("Got message:", string(message))
		if err != nil {
			fmt.Println("readPump error:", err)
			break
		}

		messageChannel <- IncomingWebsocketMessage{
			PlayerId: ws.id,
			Message:  message,
		}

		fmt.Println("Put message in channel")
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
