package handler

import (
	"net/http"
	"spt/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProjectHandler interface {
	GetProjectList(c *gin.Context)
}

type projectHandler struct {
	usecase usecase.ProjectUsecase
}

func NewProjectHandler(usecase usecase.ProjectUsecase) ProjectHandler {
	return &projectHandler{
		usecase: usecase,
	}
}

// GetProjectList  godoc
// @Summary     Get project list
// @Description Retrieves a list of projects
// @Tags        Project Management
// @Accept      application/json
// @Produce     application/json
// @Success     200 {object} map[string]interface{} "Successful response"
// @Router      /projects [get]
func (h *projectHandler) GetProjectList(c *gin.Context) {
	g := Gin{c}
	//err := h.usecase.ListProjects()
	data, err := h.usecase.GetProjectList()
	if err != nil {
		g.Response(http.StatusInternalServerError, http.StatusInternalServerError, "", err.Error())
		return
	}
	fakeData := map[string]interface{}{
		"items": []string{"testdata", ""},
	}
	fakeData["db_datas"] = data
	g.Response(http.StatusOK, http.StatusOK, "Get project list successfully", fakeData)
}
