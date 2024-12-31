package application_chatroom

import "time"

type DTO_MessagesResult struct {
	ID         string    `json:"_id"`
	ChatroomID string    `json:"chatroom_id"`
	SenderID   string    `json:"sender_id"`
	Avatar     string    `json:"avatar"`
	Nickname   string    `json:"nickname"`
	Type       string    `json:"type"`
	Data       string    `json:"data"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DTO_MembersResult struct {
	UserID   string `json:"user_id"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
}

type DTO_ChatroomResult struct {
	BaseID       string      `json:"base_id" bson:"base_id"`
	BoardID      string      `json:"board_id" bson:"board_id"`
	Title        string      `json:"title" bson:"title"`
	LastMessage  LastMessage `json:"lastMessage" bson:"lastMessage"`
	Avatar       string      `json:"avatar" bson:"avatar"`
	CreatorID    string      `json:"creator_id" bson:"creator_id"`
	CreatedAt    time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" bson:"updated_at"`
	UserList     []string    `json:"user_list" bson:"user_list"`
	Type         string      `json:"type" bson:"type"` // "private" or "board"
	ChatroomID   string      `json:"chatroom_id" bson:"chatroom_id"`
	LastReadTime time.Time   `json:"last_read_time" bson:"last_read_time"`
	UnreadCount  int         `json:"unread_count" bson:"unread_count"`
	IsMuted      bool        `json:"isMuted" bson:"is_muted"`
	IsHidden     bool        `json:"isHidden" bson:"is_hidden"`
	IsPinned     bool        `json:"isPinned" bson:"is_pinned"`
}

type LastMessage struct {
	Content      string    `json:"content" bson:"content"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	SenderID     string    `json:"senderId" bson:"senderId"`
	SenderName   string    `json:"senderName" bson:"senderName"`
	SenderAvatar string    `json:"senderAvatar" bson:"senderAvatar"`
}
