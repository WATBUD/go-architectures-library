package http

import (
	application "chat_room_mod/core/application/chatroom"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Router interface to abstract different HTTP router implementations
type ChatroomHandler struct {
	ChatroomUsecase application.ChatroomUsecase
}

func NewChatroomHandler(chatroomUsecase application.ChatroomUsecase) *ChatroomHandler {
	return &ChatroomHandler{ChatroomUsecase: chatroomUsecase}
}

func (h *ChatroomHandler) CreateChatroom(c *gin.Context) {
	var request application.DTO_CreateChatroomRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	chatroom, err := h.ChatroomUsecase.CreateChatroom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chatroom)
}

func (h *ChatroomHandler) SendMessage(c *gin.Context) {

	var input application.DTO_SendMessageRequest

	request := application.DTO_SendMessageRequest{
		ChatroomID: c.Param("chatroom_id"),
		SenderID:   c.GetHeader("Authorization"),
		Type:       input.Type,
		Data:       input.Data,
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	message, err := h.ChatroomUsecase.SendMessage(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *ChatroomHandler) GetMessages(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}

	messages, err := h.ChatroomUsecase.GetMessages(chatroomID, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *ChatroomHandler) GetChatroomMembers(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")

	memberList, err := h.ChatroomUsecase.GetChatroomMembers(chatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memberList)
}

func (h *ChatroomHandler) GetChatroomByID(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")
	chatroom, err := h.ChatroomUsecase.GetChatroomByID(chatroomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chatroom)
}

func (h *ChatroomHandler) GetChatrooms(c *gin.Context) {
	userID := c.GetHeader("Authorization")
	page := c.DefaultQuery("page", "1")          // Default page is 1
	pageSize := c.DefaultQuery("pageSize", "10") // Default page size is 10

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize parameter"})
		return
	}
	chatrooms, err := h.ChatroomUsecase.GetUserChatrooms(userID, pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chatrooms)
}

func (h *ChatroomHandler) UpdateChatroom(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")

	var input struct {
		Title  string `json:"title"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ChatroomUsecase.UpdateChatroom(chatroomID, input.Title, input.Avatar)
	if err != nil {
		if err.Error() == "chatroom not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Chatroom updated successfully"})

	//c.Status(http.StatusNoContent)
}
func (h *ChatroomHandler) LeaveChatroom(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")
	userID := c.Query("user_id")

	request := application.DTO_AddOrRemoveChatRoomUserRequest{
		ChatroomID: chatroomID,
		UserID:     userID,
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.ChatroomUsecase.LeaveChatroom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User left the chatroom successfully"})
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

func (h *ChatroomHandler) UserPinnedChatRoom(c *gin.Context) {

	var request application.DTO_PinUserRequest
	if err := c.ShouldBindJSON(&request); err != nil { // Pass a pointer to the request
		c.JSON(http.StatusBadRequest, createCustomErrorResponse(err))
		return
	}
	chatroomID := c.Param("chatroom_id")
	authHeader := c.GetHeader("Authorization")
	request.ChatroomID = chatroomID
	request.UserID = authHeader
	err := h.ChatroomUsecase.UserPinnedChatRoom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User pinned the chatroom successfully"})
}

func (h *ChatroomHandler) UserHideChatRoom(c *gin.Context) {

	var request application.DTO_HiddenUserRequest
	if err := c.ShouldBindJSON(&request); err != nil { // Pass a pointer to the request
		c.JSON(http.StatusBadRequest, createCustomErrorResponse(err))
		return
	}
	chatroomID := c.Param("chatroom_id")
	authHeader := c.GetHeader("Authorization")
	request.ChatroomID = chatroomID
	request.UserID = authHeader
	err := h.ChatroomUsecase.UserHideChatRoom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User hide the chatroom successfully"})
}

func (h *ChatroomHandler) UserMutedChatRoom(c *gin.Context) {

	var request application.DTO_MuteUserRequest
	if err := c.ShouldBindJSON(&request); err != nil { // Pass a pointer to the request
		c.JSON(http.StatusBadRequest, createCustomErrorResponse(err))
		return
	}
	chatroomID := c.Param("chatroom_id")
	authHeader := c.GetHeader("Authorization")
	request.ChatroomID = chatroomID
	request.UserID = authHeader
	err := h.ChatroomUsecase.UserMutedChatRoom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User muted the chatroom successfully"})
}

func (h *ChatroomHandler) AddUserToChatroom(c *gin.Context) {
	chatroomID := c.Param("chatroom_id")
	authHeader := c.GetHeader("Authorization")

	request := application.DTO_AddOrRemoveChatRoomUserRequest{
		ChatroomID: chatroomID,
		UserID:     authHeader,
	}
	request.ChatroomID = chatroomID
	request.UserID = authHeader

	err := h.ChatroomUsecase.AddUserToChatroom(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User joined the chatroom successfully"})
}

// func (h *ChatroomHandler) LeaveChatroom(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	chatroomID := vars["chatroom_id"]
// 	userID := r.URL.Query().Get("user_id")

// 	if chatroomID == "" || userID == "" {
// 		http.Error(w, "Chatroom ID and User ID are required", http.StatusBadRequest)
// 		return
// 	}
// 	err := h.ChatroomUsecase.RemoveUserFromChatroom(chatroomID, userID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("User left the chatroom successfully"))
// }
