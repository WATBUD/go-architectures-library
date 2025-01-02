package application

import (
	"errors"
	"folder_mod/core/domain/folder"
	"time"
)

// func (f *folder.Folder) CheckDefaultValues(operationType string) error {
// 	if f.BaseID == "" {
// 		return errors.New("BaseID is required and cannot be empty")
// 	}
// 	now := time.Now()
// 	f.UpdatedAt = now
// 	f.Type = "folder"
// 	if operationType == "create" {
// 		if f.Position == nil {
// 			// Handle case where Position is not set
// 			defaultPosition := -999.
// 			f.Position = &defaultPosition
// 		}
// 		if *f.Position < 0 {
// 			return errors.New("position must be a non-negative integer")
// 		}
// 		f.CreatedAt = now
// 	}

// 	return nil
// }

type FolderService struct {
	repo folder.FolderRepository
}

func NewFolderService(repo folder.FolderRepository) *FolderService {
	return &FolderService{repo: repo}
}

func (s *FolderService) CreateFolder(dto folder.DTO_CreateFolderRequest) (*folder.Folder, error) {
	if err := s.repo.PositionExists(dto.Base_ID, dto.Parent_ID, dto.Position); err != nil {
		return nil, err
	}

	_folder := &folder.Folder{
		BaseID:    dto.Base_ID,
		ParentID:  dto.Parent_ID,
		Position:  dto.Position,
		Data:      folder.FolderData{Color: dto.Data.Color, Name: dto.Data.Name},
		Type:      "folder",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdFolder, err := s.repo.CreateFolder(_folder)
	if err != nil {
		return nil, err
	}

	return createdFolder, nil
}
func (s *FolderService) UpdateFolderData(dto folder.DTO_UpdateFolderRequest) (*folder.Folder, error) {
	// Step 1: Retrieve the existing folder by ID
	existingFolder, err := s.repo.FindFolderByObjectID(dto.ID)
	if err != nil {
		return nil, err
	}

	// Step 2: Check if position needs to be validated
	if existingFolder.Position != dto.Position {
		if err := s.repo.PositionExists(dto.Base_ID, dto.Parent_ID, dto.Position); err != nil {
			return nil, err
		}
	}

	// Step 3: Prepare the folder for update
	_folder := &folder.Folder{
		ID:        dto.ID,
		BaseID:    dto.Base_ID,
		ParentID:  dto.Parent_ID,
		Position:  dto.Position,
		Data:      folder.FolderData{Color: dto.Data.Color, Name: dto.Data.Name},
		Type:      "folder",
		CreatedAt: existingFolder.CreatedAt, // Use the existing creation date
		UpdatedAt: time.Now(),
	}

	// Step 4: Update the folder data
	folderResult, err := s.repo.UpdateFolderData(_folder)
	if err != nil {
		return nil, err
	}

	return folderResult, nil
}

func (s *FolderService) DeleteFolder(id string) error {
	if id == "" {
		return errors.New("ID is required")
	}

	return s.repo.DeleteFolderByID(id)
}

func (s *FolderService) GetFolders() ([]*folder.Folder, error) {
	return s.repo.GetFolders()
}

func (s *FolderService) AddChildIDToParent(parentID string, childID string) error {
	return s.repo.AddChildIDToParent(parentID, childID)
}

func (s *FolderService) UpdateFolderParentID(objectID string, parentID string) error {
	return s.repo.UpdateFolderParentID(objectID, parentID)
}

func (s *FolderService) DeleteFolderByID(id string) error {
	return s.repo.DeleteFolderByID(id)
}

func (s *FolderService) PositionExists(baseID string, parentID string, position float64) error {
	return s.repo.PositionExists(baseID, parentID, position)
}

func (s *FolderService) FindFoldersByParentID(parentID string) ([]folder.Folder, error) {
	return s.repo.FindFoldersByParentID(parentID)
}
