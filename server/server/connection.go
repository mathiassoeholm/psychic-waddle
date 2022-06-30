package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type OutgoingMessage = []byte

type Connection struct {
	socket   *websocket.Conn
	playerId string
	outgoing chan OutgoingMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createConnection(writer http.ResponseWriter, request *http.Request, playerId string) (Connection, error) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return Connection{}, err
	}

	// We give the clients 60 seconds to read Ping frames.
	pongWait := 60 * time.Second

	socket.SetReadDeadline(time.Now().Add(pongWait))
	socket.SetPongHandler(func(string) error { socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	return Connection{socket: socket, outgoing: make(chan OutgoingMessage), playerId: playerId}, nil
}

func (connection *Connection) accept(messageChannel chan<- IncomingMessage, disconnect chan<- bool) {
	// TODO: how does an error in write pump stop the read pump?
	go connection.readPump(messageChannel)
	go connection.writePump(disconnect)
}

func (connection *Connection) readPump(messageChannel chan<- IncomingMessage) {
	defer connection.socket.Close()

	for {
		_, message, err := connection.socket.ReadMessage()
		fmt.Println("Got message:", string(message))
		if err != nil {
			fmt.Println("readPump error:", err)
			break
		}

		messageChannel <- IncomingMessage{
			PlayerId: connection.playerId,
			Message:  message,
		}

		fmt.Println("Put message in channel")
	}
}

func (connection *Connection) writePump(disconnect chan<- bool) {
	defer connection.socket.Close()

	for {
		select {
		case message := <-connection.outgoing:
			err := connection.socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				fmt.Println("writePump error:", err)
				disconnect <- true
				break
			}
		}
	}
}
