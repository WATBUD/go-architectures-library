package http

import (
	"folder_mod/core/application"
	"folder_mod/core/domain/folder"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FolderHandler struct {
	FolderUsecase application.FolderUsecase
}

func NewFolderHandler(folderUsecase application.FolderUsecase) *FolderHandler {
	return &FolderHandler{FolderUsecase: folderUsecase}
}

func createCustomErrorResponse(err error) gin.H {
	errorMessages := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			tag := fieldError.Tag()
			fieldName := strings.ToLower(field)
			errorMessages[fieldName] = "Field validation for " + fieldName + " failed on the '" + tag + "' tag"
		}
	}
	return gin.H{"errors": errorMessages}
}

func (h *FolderHandler) UpdateFolderData(ctx *gin.Context) {
	var request folder.DTO_UpdateFolderRequest
	id := ctx.Param("id")
	request.ID = id
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, createCustomErrorResponse(err))
		return
	}
	updatedFolder, err := h.FolderUsecase.UpdateFolderData(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Update successful",
		"folder":  updatedFolder,
	})
}

func (h *FolderHandler) CreateFolder(ctx *gin.Context) {
	var request folder.DTO_CreateFolderRequest

	// if err := ctx.ShouldBindJSON(&request); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, createCustomErrorResponse(err))
		return
	}

	insertedFolder, err := h.FolderUsecase.CreateFolder(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, insertedFolder)
}

func (h *FolderHandler) DeleteFolder(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	err := h.FolderUsecase.DeleteFolder(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Delete successful"})
}

func (h *FolderHandler) GetFolders(ctx *gin.Context) {
	folders, err := h.FolderUsecase.GetFolders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, folders)
}
