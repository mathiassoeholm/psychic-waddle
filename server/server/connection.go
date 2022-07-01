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
	playerId uint32
	outgoing chan OutgoingMessage
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func createConnection(writer http.ResponseWriter, request *http.Request, playerId uint32) (Connection, error) {
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

func (connection *Connection) accept(messageChannel chan<- IncomingMessage, disconnect chan bool) {
	go connection.readPump(messageChannel, disconnect)
	go connection.writePump(disconnect)

	go func() {
		<-disconnect
		connection.socket.Close()
	}()
}

func (connection *Connection) readPump(messageChannel chan<- IncomingMessage, disconnect chan bool) {
	for {
		select {
		case <-disconnect:
			return
		default:
			_, message, err := connection.socket.ReadMessage()
			fmt.Println("Got message:", string(message))
			if err != nil {
				fmt.Println("readPump error:", err)
				close(disconnect)
				break
			}

			messageChannel <- IncomingMessage{
				PlayerId: connection.playerId,
				Message:  message,
			}
		}
	}
}

func (connection *Connection) writePump(disconnect chan bool) {
	for {
		select {
		case <-disconnect:
			return
		case message := <-connection.outgoing:
			err := connection.socket.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				fmt.Println("writePump error:", err)
				close(disconnect)
				break
			}
		}
	}
}
