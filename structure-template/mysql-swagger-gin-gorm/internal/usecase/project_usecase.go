package usecase

import (
	"spt/internal/gorm_gen/model"
	"spt/internal/service"
)

type ProjectUsecase interface {
	ListProjects() error
	GetProjectList() ([]*model.Project, error)
}

type projectUsecase struct {
	service service.ProjectService
}

func NewProjectUsecase(service service.ProjectService) ProjectUsecase {
	return &projectUsecase{service: service}
}

func (p *projectUsecase) ListProjects() error {
	return p.service.FetchProjects()
}

func (p *projectUsecase) GetProjectList() ([]*model.Project, error) {
	return p.service.GetProjectList()
}
