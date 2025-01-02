package chatroom

type ChatroomRepository interface {
	CreateChatroom(chatroom *Entity_Chatroom) (*Entity_Chatroom, error)
	GetChatroomByID(id string) (*Entity_Chatroom, error)
	UpdateChatroom(chatroom *Entity_Chatroom) error
	SendMessage(message *Entity_Message) (*Entity_Message, error)
	GetMessages(chatroomID string, page int, pageSize int) ([]*Entity_Message, error)
	RemoveUserFromChatroom(obj SpecifyChatRoomUser) error
	UpdateChatroomUserList(obj SpecifyChatRoomUser) error

	AddUserToChatroom(obj SpecifyChatRoomUser, personal_config *ChatroomPersonalConfig) error
	SetUserChatRoomConfig(dto SetUserChatRoomConfig) error
	//GetChatroomPersonalConfig(chatroomID string) ([]*ChatroomPersonalConfig, error)
	MigrateChatroomsToUserChatrooms() error
	UpdateUserChatroomList(userID string, chatroomList []ChatroomPersonalConfig) error
	GetUserChatrooms(userID string, page, pageSize int) ([]Entity_ChatroomsList, error)
	GetUserChatroomByUserID(userID string, chatroomID string) (*ChatroomPersonalConfig, error)
	GetChatroomUserList(chatroomID string) ([]string, error)
}
