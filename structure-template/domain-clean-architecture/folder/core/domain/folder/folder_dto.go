package folder

type DTO_CreateFolderRequest struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty" binding:"required"`
	Name      string     `bson:"name" json:"name" binding:"required"`
	Base_ID   string     `bson:"base_id" json:"base_id" binding:"required"`
	Parent_ID string     `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Position  float64    `bson:"position,omitempty" json:"position,omitempty" binding:"required"`
	Data      FolderData `bson:"data" json:"data"`
}

type DTO_UpdateFolderRequest struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty" binding:"required"`
	Name      string     `bson:"name" json:"name" binding:"required"`
	Base_ID   string     `bson:"base_id" json:"base_id" binding:"required"`
	Parent_ID string     `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Position  float64    `bson:"position,omitempty" json:"position,omitempty" binding:"required"`
	Data      FolderData `bson:"data" json:"data"`
}
