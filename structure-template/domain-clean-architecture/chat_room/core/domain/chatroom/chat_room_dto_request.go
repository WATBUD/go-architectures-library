package chatroom

type DTO_CreateChatroomRequest struct {
	BaseID    string `json:"base_id" bson:"base_id" binding:"required"`
	BoardID   string `json:"board_id" bson:"board_id"`
	Title     string `json:"title" bson:"title" binding:"required"`
	CreatorID string `json:"creator_id" bson:"creator_id" binding:"required"`
	Avatar    string `json:"avatar" bson:"avatar"`
	Type      string `json:"type" bson:"type" binding:"required"`
	Nikename  string `json:"nikename" bson:"nikename"`
}

type DTO_SendMessageRequest struct {
	ChatroomID string `json:"chatroom_id" binding:"required"`
	SenderID   string `json:"sender_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Data       string `json:"data" binding:"required"`
}

type DTO_LeaveChatroomRequest struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" binding:"required"`
	UserID     string `json:"user_id" bson:"user_id" binding:"required"`
}

type DTO_AddOrRemoveChatRoomUserRequest struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" binding:"required"`
	UserID     string `json:"user_id" bson:"user_id" binding:"required"`
}

type DTO_SetUserChatRoomConfig struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" `
	UserID     string `json:"user_id" bson:"user_id"`
	ParamName  string `json:"param_name" bson:"param_name"`
	ParamValue any    `json:"param_value" bson:"param_value"`
}

type DTO_MuteUserRequest struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" `
	UserID     string `json:"user_id" bson:"user_id"`
	Mute       *bool  `json:"mute" bson:"mute" binding:"required"`
}

type DTO_PinUserRequest struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" `
	UserID     string `json:"user_id" bson:"user_id"`
	Pin        *bool  `json:"pin" bson:"pin" binding:"required"`
}
