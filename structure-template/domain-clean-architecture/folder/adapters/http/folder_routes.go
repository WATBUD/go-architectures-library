package http

import (
	"github.com/gin-gonic/gin"
)

func SetupChatroomRoutes(router *gin.Engine, folderHandler *FolderHandler) {
	router.POST("/folders", folderHandler.CreateFolder)
	router.GET("/folders", folderHandler.GetFolders)
	router.DELETE("/folders/:id", folderHandler.DeleteFolder)
	router.PUT("/folders/:id", folderHandler.UpdateFolderData)

}
