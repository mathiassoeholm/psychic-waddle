package server

const (
	newConnection   = "new_connection"
	receivedMessage = "received_message"
)

type ConnectionEvent interface {
	Type() string
}

type NewConnection struct {
	PlayerId uint32
}

func (NewConnection) Type() string {
	return newConnection
}

type ReceivedMessage struct {
	PlayerId uint32
	Message  []byte
}

func (ReceivedMessage) Type() string {
	return receivedMessage
}
