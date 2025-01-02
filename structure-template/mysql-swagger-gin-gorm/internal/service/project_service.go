package service

import (
	"context"
	"spt/internal/gorm_gen/model"
	"spt/internal/repository"
)

type ProjectService interface {
	FetchProjects() error
	GetProjectList() ([]*model.Project, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{
		repo: repo,
	}
}

func (p *projectService) GetProjectList() ([]*model.Project, error) {
	return p.repo.GetProjectList(context.TODO())
}

func (p *projectService) FetchProjects() error {
	return p.repo.FindAll()
}
