package connections

type WebsocketEventType = string

const (
	WebsocketEventNewConnectionType WebsocketEventType = "new_connection"
)

type WebsocketEvent interface {
	Type() WebsocketEventType
}
