package router

import (
	"github.com/gin-gonic/gin"
	"spt/internal/db"
	"spt/internal/handler"
	"spt/internal/repository"
	"spt/internal/service"
	"spt/internal/usecase"
)

func RegisterProjectRoutes(r *gin.RouterGroup) {
	db := db.Instance()
	projectRepository := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepository)
	projectUsecase := usecase.NewProjectUsecase(projectService)
	projectHandlers := handler.NewProjectHandler(projectUsecase)

	r.GET("/projects", projectHandlers.GetProjectList)

}
