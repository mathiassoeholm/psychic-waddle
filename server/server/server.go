package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type IncomingMessage struct {
	PlayerId uint32
	Message  []byte
}

type Server struct {
	port                 int
	events               chan ConnectionEvent
	nextId               uint32
	playerIdToConnection map[uint32]Connection

	incomingMessages chan IncomingMessage
}

func New(port int) *Server {
	return &Server{
		port:                 port,
		events:               make(chan ConnectionEvent, 250),
		playerIdToConnection: make(map[uint32]Connection),
		incomingMessages:     make(chan IncomingMessage),
	}
}

func (server *Server) Events() <-chan ConnectionEvent {
	return server.events
}

func (server *Server) Run() {
	fmt.Println("Server running on port", server.port)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		playerId := server.nextId
		server.nextId++

		connection, err := createConnection(writer, request, playerId)

		if err != nil {
			fmt.Println("Error creating connection:", err)
			return
		}

		server.playerIdToConnection[playerId] = connection

		disconnect := make(chan bool)
		connection.accept(server.incomingMessages, disconnect)

		go func() {
			<-disconnect
			delete(server.playerIdToConnection, playerId)
		}()

		server.events <- NewConnection{
			PlayerId: playerId,
		}

		fmt.Println("New connection:", playerId)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(server.port), serveMux)
		if err != nil {
			log.Fatal("Websocket server error:", err)
		}
	}()

	for {
		select {
		case message, ok := <-server.incomingMessages:
			if !ok {
				panic("incoming messages channel closed")
			}

			fmt.Printf("Incoming message from %q: %v\n", message.PlayerId, string(message.Message))

			if len(message.Message) > 0 {
				server.handleMessage(message.PlayerId, message.Message[0], message.Message[1:])
			}

			server.events <- ReceivedMessage{
				PlayerId: message.PlayerId,
				Message:  message.Message,
			}
		}
	}
}

func (server *Server) SendToAllExcept(playerId uint32, message []byte) {
	for _, connection := range server.playerIdToConnection {
		if connection.playerId != playerId {
			fmt.Printf("Sending message to %q", connection.playerId)
			connection.outgoing <- message
		}
	}
}

func (server *Server) Send(playerId uint32, message []byte) error {
	connection, exists := server.playerIdToConnection[playerId]
	if !exists {
		return fmt.Errorf("no connection with player id %q", playerId)
	}
	connection.outgoing <- message
	return nil
}
