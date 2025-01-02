package chatroom

type ChatroomService interface {
	CreateChatroom(dto DTO_CreateChatroomRequest) (*Entity_Chatroom, error)
	GetChatroomByID(id string) (*Entity_Chatroom, error)
	GetChatrooms(userID string) ([]*Entity_Chatroom, error)
	UpdateChatroom(chatroomID string, Title string, Avatar string) error
	SendMessage(dto DTO_SendMessageRequest) (*Entity_Message, error)
	GetMessages(chatroomID string) ([]*Entity_Message, error)
	RemoveUserFromChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
}
