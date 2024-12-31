package chatroom

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
