package folder

import (
	"time"
)

// Folder represents a folder in MongoDB
type Folder struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty"`
	BaseID    string     `bson:"base_id" json:"base_id" binding:"required"`
	ParentID  string     `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Position  float64    `bson:"position,omitempty" json:"position,omitempty"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	Data      FolderData `bson:"data" json:"data"`
	Type      string     `bson:"type" json:"type"`
}

// FolderData represents custom data stored in a folder
type FolderData struct {
	Name  string `bson:"name" json:"name"`
	Color string `bson:"color" json:"color"`
}

// Board represents a board in MongoDB
type Board struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	BaseID    string    `bson:"base_id" json:"base_id"`
	Position  float64   `bson:"position" json:"position"`
	ParentID  string    `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	Tabs      []Tab     `bson:"tabs" json:"tabs"`
	Data      BoardData `bson:"data" json:"data"`
	Type      string    `bson:"type" json:"type"`
}

// BoardData represents custom data stored in a board
type BoardData struct {
	Name  string `bson:"name" json:"name"`
	Color string `bson:"color" json:"color"`
}

// Tab represents a tab in a board
type Tab struct {
	Permission Permission `bson:"permission" json:"permission"`
}

// Permission represents permissions for a tab
type Permission struct {
	Read     bool   `bson:"read" json:"read"`
	UserList []User `bson:"userlist" json:"userlist"`
}

// User represents a user with specific permissions
type User struct {
	UserID string `bson:"userid" json:"userid"`
}
