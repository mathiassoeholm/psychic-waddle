package main

import (
	"fmt"

	server "github.com/mathiassoeholm/psychic-waddle/server/server"
)

func main() {
	s := server.New(4000)
	go s.Run()
	for event := range s.Events() {
		switch casted := event.(type) {
		case server.NewConnection:
			s.Send(casted.PlayerId, []byte("Welcome to the game! :-)"))
		case server.ReceivedMessage:
			fmt.Println("Got message:", string(casted.Message))
			s.SendToAllExcept(casted.PlayerId, casted.Message)
		}
	}
}
