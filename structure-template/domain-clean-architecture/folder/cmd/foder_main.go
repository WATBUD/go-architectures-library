package main

import (
	"context"
	adapters_http "folder_mod/adapters/http"
	"folder_mod/adapters/repository"
	"folder_mod/core/application"
	mongo_database "folder_mod/infrastructure"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	client, err := mongo_database.NewMongoClient(ctx)
	if err != nil {
		panic(err)
	}

	folderRepo := repository.NewMongoFolderRepository(client)
	folderService := application.NewFolderService(folderRepo)
	folderUsecase := application.NewFolderUsecase(folderService)
	folderHandler := adapters_http.NewFolderHandler(folderUsecase)
	router := gin.Default()
	adapters_http.SetupChatroomRoutes(router, folderHandler)
	router.Run(":8080")

}
