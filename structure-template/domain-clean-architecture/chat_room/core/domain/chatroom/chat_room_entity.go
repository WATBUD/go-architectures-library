package chatroom

import (
	"time"
)

type UserChatroomConfig struct {
	ChatroomID string
	UserID     string
	ParamName  string
	ParamValue interface{}
}
type LastMessage struct {
	Content      string    `json:"content" bson:"content"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	SenderID     string    `json:"senderId" bson:"senderId"`
	SenderName   string    `json:"senderName" bson:"senderName"`
	SenderAvatar string    `json:"senderAvatar" bson:"senderAvatar"`
}

type Entity_Chatroom struct {
	ID            string    `json:"_id" bson:"_id"`
	BaseID        string    `json:"base_id" bson:"base_id"`
	BoardID       string    `json:"board_id" bson:"board_id"`
	Title         string    `json:"title" bson:"title"`
	LastMessageID string    `json:"last_message_id" bson:"last_message_id"`
	Avatar        string    `json:"avatar" bson:"avatar"`
	CreatorID     string    `json:"creator_id" bson:"creator_id"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
	UserList      []string  `json:"user_list" bson:"user_list"`
	Type          string    `json:"type" bson:"type"` // "private" or "board"
}

type Entity_ChatroomsList struct {
	Chatroom       Entity_Chatroom
	PersonalConfig ChatroomPersonalConfig
	LastMessage    LastMessage
}

type ChatroomPersonalConfig struct {
	ChatroomID   string    `json:"chatroom_id" bson:"chatroom_id"`
	LastReadTime time.Time `json:"last_read_time" bson:"last_read_time"`
	UnreadCount  int       `json:"unread_count" bson:"unread_count"`
	IsMuted      bool      `json:"isMuted" bson:"is_muted"`
	IsHidden     bool      `json:"isHidden" bson:"is_hidden"`
	IsPinned     bool      `json:"isPinned" bson:"is_pinned"`
}

type Entity_UserChatroom struct {
	UserID       string                   `json:"user_id" bson:"user_id"`
	ChatroomList []ChatroomPersonalConfig `json:"chatroom_list" bson:"chatroom_list"`
}

type Entity_Message struct {
	ID         string    `json:"_id" bson:"_id"`
	ChatroomID string    `json:"chatroom_id" bson:"chatroom_id"`
	SenderID   string    `json:"sender_id" bson:"sender_id"`
	Type       string    `json:"type" bson:"type"`
	Data       string    `json:"data" bson:"data"`
	Avatar     string    `json:"avatar" bson:"avatar"`
	Nickname   string    `json:"nickname" bson:"nickname"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

type SetUserChatRoomConfig struct {
	ChatroomID string `json:"chatroom_id" bson:"chatroom_id" `
	UserID     string `json:"user_id" bson:"user_id"`
	ParamName  string `json:"param_name" bson:"param_name"`
	ParamValue any    `json:"param_value" bson:"param_value"`
}
