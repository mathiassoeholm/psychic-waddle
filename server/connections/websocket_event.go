package connections

type WebsocketEventType = string

const (
	WebsocketEventNewConnectionType WebsocketEventType = "new_connection"
	WebsocketEventMessageType       WebsocketEventType = "message"
)

type WebsocketEvent interface {
	Type() WebsocketEventType
}

type NewConnectionWebsocketEvent struct {
	PlayerId string
}

func (NewConnectionWebsocketEvent) Type() WebsocketEventType {
	return WebsocketEventNewConnectionType
}

var _ WebsocketEvent = NewConnectionWebsocketEvent{}

type MessageWebsocketEvent struct {
	PlayerId string
	Message  []byte
}

func (MessageWebsocketEvent) Type() WebsocketEventType {
	return WebsocketEventMessageType
}

var _ WebsocketEvent = MessageWebsocketEvent{}
