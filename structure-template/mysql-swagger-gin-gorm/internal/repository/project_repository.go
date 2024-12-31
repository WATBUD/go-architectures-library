package repository

import (
	"context"
	"errors"
	"spt/internal/db"
	"spt/internal/gorm_gen/model"
)

type ProjectRepository interface {
	FindAll() error
	GetProjectList(ctx context.Context) ([]*model.Project, error)
}

type projectRepository struct {
	dbService db.Service
}

func NewProjectRepository(dbService db.Service) ProjectRepository {
	return &projectRepository{
		dbService: dbService,
	}
}

func (p *projectRepository) GetProjectList(ctx context.Context) ([]*model.Project, error) {
	var projectsFromDB []*model.Project
	tx := p.dbService.GetDB().Find(&projectsFromDB)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var projects []*model.Project
	for _, project := range projectsFromDB {
		projects = append(projects, &model.Project{
			ID:          project.ID,
			ProjectName: project.ProjectName,
		})
	}
	return projects, nil
}

func (p *projectRepository) FindAll() error {
	return errors.New("not implemented")
}
