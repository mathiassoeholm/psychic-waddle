package connections

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type WebsocketServer struct {
	port                 int
	events               chan WebsocketEvent
	nextId               int
	playerIdToConnection map[string]Websocket

	incoming_connections chan Websocket
}

func NewWebsocketServer(port int) *WebsocketServer {
	return &WebsocketServer{
		port:                 port,
		events:               make(chan WebsocketEvent, 250),
		playerIdToConnection: make(map[string]Websocket),
		incoming_connections: make(chan Websocket),
	}
}

func (ws *WebsocketServer) Events() <-chan WebsocketEvent {
	return ws.events
}

func (ws *WebsocketServer) Run() {
	fmt.Println("Websocket server running on port", ws.port)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		serveWebsocket(ws.incoming_connections, writer, request)
	})

	go func() {
		err := http.ListenAndServe(":"+strconv.Itoa(ws.port), serveMux)
		if err != nil {
			log.Fatal("Websocket server error:", err)
		}
	}()

	for {
		select {
		case connection := <-ws.incoming_connections:
			playerId := strconv.Itoa(ws.nextId)
			ws.events <- NewConnectionWebsocketEvent{
				PlayerId: playerId,
			}
			ws.playerIdToConnection[playerId] = connection
		}
	}
}
