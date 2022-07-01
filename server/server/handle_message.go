package server

func (server *Server) handleMessage(playerId uint32, messageId byte, message []byte) {
	switch messageId {

	case sendChatMessageId:
		server.SendToAllExcept(playerId, append(append([]byte{receiveChatMessageId}, encodeUint32(playerId)...), message...))
	}
}
