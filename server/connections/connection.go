package connections

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type OutgoingMessage = []byte

type Connection struct {
	socket   *websocket.Conn
	id       string
	outgoing chan OutgoingMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createConnection(writer http.ResponseWriter, request *http.Request) (Connection, error) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return Connection{}, err
	}

	// We give the clients 60 seconds to read Ping frames.
	pongWait := 60 * time.Second

	socket.SetReadDeadline(time.Now().Add(pongWait))
	socket.SetPongHandler(func(string) error { socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	return Connection{socket: socket, outgoing: make(chan OutgoingMessage)}, nil
}

func (ws *Connection) accept(playerId string, messageChannel chan<- IncomingMessage) {
	ws.id = playerId

	go ws.readPump(messageChannel)
	go ws.writePump()

}

func (ws *Connection) readPump(messageChannel chan<- IncomingMessage) {
	defer ws.socket.Close()

	for {
		_, message, err := ws.socket.ReadMessage()
		fmt.Println("Got message:", string(message))
		if err != nil {
			fmt.Println("readPump error:", err)
			break
		}

		messageChannel <- IncomingMessage{
			PlayerId: ws.id,
			Message:  message,
		}

		fmt.Println("Put message in channel")
	}
}

func (ws *Connection) writePump() {
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
