package folder_mod

import (
	"context"
	"folder_mod/adapters/repository"
	"folder_mod/core/application"
	"folder_mod/core/domain/folder"
	"sync"
	"time"

	goa "goa.design/goa/v3/pkg"
	goa_folder "mai.today/api/gen/folder"
	"mai.today/authentication"
	"mai.today/database/mongodb"
	"mai.today/realtime"
)

func NewFolderServices(repo folder.FolderRepository) *FolderService {
	return &FolderService{
		applicationUsecase: application.NewFolderUsecase(application.NewFolderService(repo)),
	}
}

func (fs *FolderService) CreateFolder(ctx context.Context, p *goa_folder.CreateFolderPayload) (*goa_folder.CreateFolderResult, error) {
	dto := folder.DTO_CreateFolderRequest{
		Base_ID:   p.BaseID,
		Parent_ID: p.ParentID,
		Position:  p.Position,
		Data:      folder.FolderData{Color: p.Data.Color, Name: p.Data.Name},
	}
	newFolder, err := fs.applicationUsecase.CreateFolder(dto)
	if err != nil {
		return nil, err
	}
	return &goa_folder.CreateFolderResult{
		Command:   &goa_folder.Command{Type: "CreateFolder"},
		Timestamp: newFolder.CreatedAt.Unix(),
		Data: &goa_folder.Createfolderresultdata{
			FolderID:  newFolder.ID,
			BaseID:    newFolder.BaseID,
			ParentID:  newFolder.ParentID,
			Position:  newFolder.Position,
			CreatedAt: newFolder.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: newFolder.UpdatedAt.UTC().Format(time.RFC3339),
			Data:      &goa_folder.FolderData{Color: newFolder.Data.Color, Name: newFolder.Data.Name},
			Type:      newFolder.Type,
		},
	}, nil
}

func (fs *FolderService) DeleteFolder(ctx context.Context, p *goa_folder.DeleteFolderPayload) (*goa_folder.DeleteFolderResult, error) {
	err := fs.applicationUsecase.DeleteFolderByID(p.FolderID)
	if err != nil {
		return nil, err
	}
	return &goa_folder.DeleteFolderResult{
		Command:   &goa_folder.Command{Type: "DeleteFolder"},
		Timestamp: time.Now().Unix(),
		Data: &goa_folder.Createfolderresultdata{
			FolderID: p.FolderID,
		},
	}, nil
}

func (fs *FolderService) UpdateFolder(ctx context.Context, p *goa_folder.UpdateFolderPayload) (*goa_folder.UpdateFolderResult, error) {
	dto := folder.DTO_UpdateFolderRequest{
		ID:        p.FolderID,
		Base_ID:   p.BaseID,
		Parent_ID: p.ParentID,
		Position:  p.Position,
		Data:      folder.FolderData{Color: p.Data.Color, Name: p.Data.Name},
	}

	_folder, err := fs.applicationUsecase.UpdateFolderData(dto)
	if err != nil {
		return nil, err
	}

	return &goa_folder.UpdateFolderResult{
		Command:   &goa_folder.Command{Type: "UpdateFolder"},
		Timestamp: _folder.UpdatedAt.Unix(),
		Data: &goa_folder.Createfolderresultdata{
			FolderID:  _folder.ID,
			BaseID:    _folder.BaseID,
			ParentID:  _folder.ParentID,
			Position:  _folder.Position, // No need to use pointer
			CreatedAt: _folder.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: _folder.UpdatedAt.UTC().Format(time.RFC3339),
			Data:      &goa_folder.FolderData{Color: _folder.Data.Color, Name: _folder.Data.Name}, // No pointer
			Type:      _folder.Type,
		},
	}, nil
}

func (fs *FolderService) ListFolders(ctx context.Context, p *goa_folder.ListFoldersPayload) (*goa_folder.ListFoldersResult, error) {
	folders, err := fs.applicationUsecase.GetFolders()
	if err != nil {
		return nil, err
	}
	var folderCollection []*goa_folder.Createfolderresultdata
	for _, f := range folders {
		folderCollection = append(folderCollection, &goa_folder.Createfolderresultdata{
			FolderID:  f.ID,
			BaseID:    f.BaseID,
			ParentID:  f.ParentID,
			Position:  f.Position,
			CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: f.UpdatedAt.UTC().Format(time.RFC3339),
			Data:      &goa_folder.FolderData{Color: f.Data.Color, Name: f.Data.Name},
			Type:      f.Type,
		})
	}
	result := &goa_folder.ListFoldersResult{
		Command:   &goa_folder.Command{Type: "GetListFolders"},
		Timestamp: time.Now().Unix(),
		Data:      folderCollection,
	}
	return result, nil
}

// MakeInvalidToken builds a goa.ServiceError from an error.
func MakeInvalidToken(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "invalid token", false, false, false)
}

func (fs *FolderService) ReceiveCreateFolder(context.Context, *goa_folder.ReceiveCreateFolderPayload) (res *goa_folder.CreateFolderResult, err error) {

	panic("not implemented")
}
func (fs *FolderService) ReceiveDeleteFolder(context.Context, *goa_folder.ReceiveDeleteFolderPayload) (res *goa_folder.DeleteFolderResult, err error) {
	panic("not implemented")
}
func (fs *FolderService) ReceiveUpdateFolder(context.Context, *goa_folder.ReceiveUpdateFolderPayload) (res *goa_folder.UpdateFolderResult, err error) {
	panic("not implemented")

}
func (fs *FolderService) ReceiveListFolders(context.Context, *goa_folder.ReceiveListFoldersPayload) (res *goa_folder.ListFoldersResult, err error) {
	panic("not implemented")
}

var (
	once     sync.Once
	instance *FolderService
)

func Instance() *FolderService {
	once.Do(func() {
		instance = newFolderService()
	})
	return instance
}

type FolderService struct {
	realtime realtime.Realtime
	goa_folder.Auther
	applicationUsecase application.FolderUsecase
	applicationService folder.FolderService
}

func newFolderService() *FolderService {
	applicationService := application.NewFolderService(repository.NewMongoFolderRepository(mongodb.Instance()))
	applicationUsecase := application.NewFolderUsecase(applicationService)
	return &FolderService{
		realtime:           realtime.Instance(),
		Auther:             authentication.Instance(),
		applicationService: applicationService,
		applicationUsecase: applicationUsecase,
	}
}
