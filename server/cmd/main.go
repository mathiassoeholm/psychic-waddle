package main

import "github.com/mathiassoeholm/psychic-waddle/server/connections"

func main() {
	server := connections.NewWebsocketServer(4000)
	server.Run()
}
