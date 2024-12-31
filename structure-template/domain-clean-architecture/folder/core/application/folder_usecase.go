package application

import (
	"folder_mod/core/domain/folder"
)

type FolderUsecase interface {
	CreateFolder(dto folder.DTO_CreateFolderRequest) (*folder.Folder, error)
	DeleteFolder(id string) error
	GetFolders() ([]*folder.Folder, error)
	UpdateFolderData(dto folder.DTO_UpdateFolderRequest) (*folder.Folder, error)
	DeleteFolderByID(id string) error
	UpdateFolderParentID(baseID string, parentID string) error
	AddChildIDToParent(parentID string, childID string) error
	PositionExists(baseID string, parentID string, position float64) error
	FindFoldersByParentID(parentID string) ([]folder.Folder, error)
}

type folderUsecase struct {
	folderService folder.FolderService
}

func NewFolderUsecase(folderService folder.FolderService) FolderUsecase {
	return &folderUsecase{folderService: folderService}
}

func (u *folderUsecase) CreateFolder(dto folder.DTO_CreateFolderRequest) (*folder.Folder, error) {
	return u.folderService.CreateFolder(dto)
}

func (u *folderUsecase) DeleteFolder(id string) error {
	return u.folderService.DeleteFolderByID(id) // 调用 DeleteFolderByID
}

func (u *folderUsecase) GetFolders() ([]*folder.Folder, error) {
	return u.folderService.GetFolders()
}

func (u *folderUsecase) UpdateFolderData(dto folder.DTO_UpdateFolderRequest) (*folder.Folder, error) {
	return u.folderService.UpdateFolderData(dto)
}

func (u *folderUsecase) DeleteFolderByID(id string) error {
	return u.folderService.DeleteFolderByID(id)
}

func (u *folderUsecase) UpdateFolderParentID(baseID string, parentID string) error {
	return u.folderService.UpdateFolderParentID(baseID, parentID)
}

func (u *folderUsecase) AddChildIDToParent(parentID string, childID string) error {
	return u.folderService.AddChildIDToParent(parentID, childID)
}

func (u *folderUsecase) PositionExists(baseID string, parentID string, position float64) error {
	return u.folderService.PositionExists(baseID, parentID, position)
}

func (u *folderUsecase) FindFoldersByParentID(parentID string) ([]folder.Folder, error) {
	return u.folderService.FindFoldersByParentID(parentID)
}
