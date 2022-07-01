package server

const (
	sendChatMessageId byte = iota
	receiveChatMessageId
	positionUpdateMessageId
)

func encodeUint32(value uint32) []byte {
	return []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value)}
}
