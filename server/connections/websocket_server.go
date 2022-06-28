package connections

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

type WebsocketServer struct {
	port                 int
	events               chan ConnectionEvent
	nextId               int
	playerIdToConnection map[string]Connection

	incomingConnections chan Connection
	incomingMessages    chan IncomingMessage
}

func NewWebsocketServer(port int) *WebsocketServer {
	return &WebsocketServer{
		port:                 port,
		events:               make(chan ConnectionEvent, 250),
		playerIdToConnection: make(map[string]Connection),
		incomingConnections:  make(chan Connection),
		incomingMessages:     make(chan IncomingMessage),
	}
}

func (ws *WebsocketServer) Events() <-chan ConnectionEvent {
	return ws.events
}

func (ws *WebsocketServer) Run() {
	fmt.Println("Websocket server running on port", ws.port)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		connection, err := createConnection(writer, request)
		if err != nil {
			fmt.Println("Error creating connection:", err)
			return
		}

		ws.incomingConnections <- connection
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(ws.port), serveMux)
		if err != nil {
			log.Fatal("Websocket server error:", err)
		}
	}()

	for {
		select {
		case connection := <-ws.incomingConnections:
			playerId := strconv.Itoa(ws.nextId)
			ws.nextId++

			ws.playerIdToConnection[playerId] = connection
			connection.accept(playerId, ws.incomingMessages)

			ws.events <- NewConnection{
				PlayerId: playerId,
			}

			fmt.Println("New connection:", playerId)

		case message, ok := <-ws.incomingMessages:
			if !ok {
				panic("incoming_messages channel closed")
			}

			fmt.Printf("Incoming message from %q: %v\n", message.PlayerId, string(message.Message))
		}
	}
}

func (ws *WebsocketServer) Send(playerId string, message []byte) error {
	connection, exists := ws.playerIdToConnection[playerId]
	if !exists {
		return fmt.Errorf("no connection with player id %q", playerId)
	}
	connection.outgoing <- message
	return nil
}
