package connections

type WebsocketEventType = string

const (
	WebsocketEventNewConnectionType WebsocketEventType = "new_connection"
)

type WebsocketEvent interface {
	Type() WebsocketEventType
}

type NewConnectionWebsocketEvent struct {
	PlayerId string
}

// Type implements WebsocketEvent
func (NewConnectionWebsocketEvent) Type() string {
	return WebsocketEventNewConnectionType
}

var _ WebsocketEvent = NewConnectionWebsocketEvent{}
