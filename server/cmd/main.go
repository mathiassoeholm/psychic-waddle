package main

import "github.com/mathiassoeholm/psychic-waddle/server/connections"

func main() {
	server := connections.NewWebsocketServer(4000)
	go server.Run()
	for event := range server.Events() {
		switch casted := event.(type) {
		case connections.NewConnection:
			server.Send(casted.PlayerId, []byte("Welcome to the game! :-)"))
		}
	}
}
