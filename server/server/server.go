package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type IncomingMessage struct {
	PlayerId string
	Message  []byte
}

type Server struct {
	port                 int
	events               chan ConnectionEvent
	nextId               int
	playerIdToConnection map[string]Connection

	incomingConnections chan Connection
	incomingMessages    chan IncomingMessage
}

func New(port int) *Server {
	return &Server{
		port:                 port,
		events:               make(chan ConnectionEvent, 250),
		playerIdToConnection: make(map[string]Connection),
		incomingConnections:  make(chan Connection),
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
		connection, err := createConnection(writer, request)
		if err != nil {
			fmt.Println("Error creating connection:", err)
			return
		}

		server.incomingConnections <- connection
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(server.port), serveMux)
		if err != nil {
			log.Fatal("Websocket server error:", err)
		}
	}()

	for {
		select {
		case connection := <-server.incomingConnections:
			playerId := strconv.Itoa(server.nextId)
			server.nextId++

			server.playerIdToConnection[playerId] = connection
			connection.accept(playerId, server.incomingMessages)

			server.events <- NewConnection{
				PlayerId: playerId,
			}

			fmt.Println("New connection:", playerId)

		case message, ok := <-server.incomingMessages:
			if !ok {
				panic("incoming messages channel closed")
			}

			fmt.Printf("Incoming message from %q: %v\n", message.PlayerId, string(message.Message))
		}
	}
}

func (server *Server) Send(playerId string, message []byte) error {
	connection, exists := server.playerIdToConnection[playerId]
	if !exists {
		return fmt.Errorf("no connection with player id %q", playerId)
	}
	connection.outgoing <- message
	return nil
}
