package folder

type FolderService interface {
	CreateFolder(dto DTO_CreateFolderRequest) (*Folder, error)
	GetFolders() ([]*Folder, error)
	UpdateFolderData(dto DTO_UpdateFolderRequest) (*Folder, error)
	DeleteFolderByID(id string) error
	UpdateFolderParentID(baseID string, parentID string) error
	AddChildIDToParent(parentID string, childID string) error
	PositionExists(baseID string, parentID string, position float64) error
	FindFoldersByParentID(parentID string) ([]Folder, error)
}
