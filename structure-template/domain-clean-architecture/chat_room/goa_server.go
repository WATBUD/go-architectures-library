package chat_room_mod

import (
	"goa.design/goa/v3/http"
	"mai.today/api/gen/chatroom"
	"mai.today/api/gen/http/chatroom/server"
)

func NewChatRoomServer() *server.Server {
	endpoints := chatroom.NewEndpoints(Instance())

	dec := http.RequestDecoder
	enc := http.ResponseEncoder
	svr := server.New(endpoints, http.NewMuxer(), dec, enc, nil, nil)

	return svr
}
