package http

import (
	"github.com/gin-gonic/gin"
)

func SetupChatroomRoutes(router *gin.Engine, chatroomHandler *ChatroomHandler) {
	router.GET("/chatroom/check/:chatroom_id", chatroomHandler.GetChatroomByID)
	router.GET("/chatroom/room/", chatroomHandler.GetChatrooms)
	router.GET("/chatroom/:chatroom_id/messages", chatroomHandler.GetMessages)
	router.GET("/chatroom/:chatroom_id/members", chatroomHandler.GetChatroomMembers)

	router.POST("/chatroom", chatroomHandler.CreateChatroom)
	router.POST("/chatroom/:chatroom_id/add-user", chatroomHandler.AddUserToChatroom)
	router.POST("/chatroom/:chatroom_id/messages", chatroomHandler.SendMessage)

	router.PUT("/chatroom/:chatroom_id", chatroomHandler.UpdateChatroom)
	router.PUT("/chatroom/:chatroom_id/user-muted-chatroom", chatroomHandler.UserMutedChatRoom)
	router.PUT("/chatroom/:chatroom_id/user-pinned-chatroom", chatroomHandler.UserPinnedChatRoom)
	router.PUT("/chatroom/:chatroom_id/user-hide-chatroom", chatroomHandler.UserHideChatRoom)

	router.DELETE("/chatroom/:chatroom_id/leave-chatroom", chatroomHandler.LeaveChatroom)
}
