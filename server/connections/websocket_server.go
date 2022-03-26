package connections

import "fmt"

type WebsocketServer struct {
	port    int
	channel chan WebsocketEvent
}

func NewWebSocketServer(port int) *WebsocketServer {
	return &WebsocketServer{
		port: port, channel: make(chan WebsocketEvent),
	}
}

func (ws *WebsocketServer) Events() <-chan WebsocketEvent {
	return ws.channel
}

func (ws *WebsocketServer) Run() {
	fmt.Println("Websocket server running on port", ws.port)
}
