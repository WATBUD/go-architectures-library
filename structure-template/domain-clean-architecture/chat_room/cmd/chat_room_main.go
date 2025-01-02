package main

import (
	adapters_http "chat_room_mod/adapters/http"
	"chat_room_mod/adapters/repository"
	application "chat_room_mod/core/application/chatroom"
	mongo_database "chat_room_mod/infrastructure"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()
	client, err := mongo_database.NewMongoClient(ctx)
	if err != nil {
		panic(err)
	}

	chatroomRepo := repository.NewMongoChatroomRepository(client)
	chatroomService := application.NewChatroomService(chatroomRepo)
	chatroomUsecase := application.NewChatroomUsecase(chatroomService)
	chatroomHandler := adapters_http.NewChatroomHandler(chatroomUsecase)
	router := gin.Default()
	adapters_http.SetupChatroomRoutes(router, chatroomHandler)
	// Start the server on port 8080
	router.Run(":8080")
}
