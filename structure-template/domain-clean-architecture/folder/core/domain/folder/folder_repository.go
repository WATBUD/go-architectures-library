package folder

type FolderRepository interface {
	CreateFolder(folder *Folder) (*Folder, error)
	GetFolders() ([]*Folder, error)
	UpdateFolderData(folder *Folder) (*Folder, error)
	DeleteFolderByID(id string) error
	UpdateFolderParentID(baseID string, parentID string) error
	AddChildIDToParent(parentID string, childID string) error
	PositionExists(baseID string, parentID string, position float64) error
	FindFoldersByParentID(parentID string) ([]Folder, error)
	FindFolderByObjectID(objectID string) (*Folder, error)
}
