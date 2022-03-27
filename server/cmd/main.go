package main

import "github.com/mathiassoeholm/psychic-waddle/server/connections"

func main() {
	server := connections.NewWebsocketServer(4000)
	go server.Run()
	for {
		select {
		case event := <-server.Events():
			switch casted := event.(type) {
			case connections.NewConnectionWebsocketEvent:
				server.Send(casted.PlayerId, []byte("Welcome to the game!"))
			}
		}
	}
}
