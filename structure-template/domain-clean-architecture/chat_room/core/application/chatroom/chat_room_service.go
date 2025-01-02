package application_chatroom

import (
	//application "chat_room_mod/core/application/application_chatroom"
	domain "chat_room_mod/core/domain/chatroom"
	"errors"
	"fmt"
	"time"
)

type ChatroomService interface {
	CreateChatroom(chatroom DTO_CreateChatroomRequest) (*domain.Entity_Chatroom, error)
	GetChatroomByID(id string) (*domain.Entity_Chatroom, error)
	UpdateChatroom(chatroomID string, Title string, Avatar string) error
	SendMessage(dto DTO_SendMessageRequest) (*DTO_MessagesResult, error)
	GetMessages(chatroomID string, page int, pageSize int) ([]*domain.Entity_Message, error)
	RemoveUserFromChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	LeaveChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	AddUserToChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	SetUserChatRoomConfig(dto DTO_SetUserChatRoomConfig) error
	GetChatroomUserList(chatroomID string) ([]string, error)
	GetUserChatrooms(userID string, page int, pageSize int) ([]DTO_ChatroomResult, error)
}

type chatroomService struct {
	repo domain.ChatroomRepository
}

func NewChatroomService(repo domain.ChatroomRepository) *chatroomService {
	return &chatroomService{repo: repo}
}

func (s *chatroomService) SetUserChatRoomConfig(dto DTO_SetUserChatRoomConfig) error {
	chatroomData, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil || chatroomData == nil {
		return errors.New("chatroom does not exist")
	}

	_SetUserChatRoomConfig := domain.SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  dto.ParamName,
		ParamValue: dto.ParamValue,
	}

	return s.repo.SetUserChatRoomConfig(_SetUserChatRoomConfig)
}

func (s *chatroomService) CreateChatroom(dto DTO_CreateChatroomRequest) (*domain.Entity_Chatroom, error) {

	chatroom := &domain.Entity_Chatroom{
		Title:         dto.Title,
		CreatorID:     dto.CreatorID,
		BaseID:        dto.BaseID,
		BoardID:       dto.BoardID,
		Avatar:        dto.Avatar,
		LastMessageID: "",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Type:          "private",
		UserList:      []string{dto.CreatorID},
	}

	return s.repo.CreateChatroom(chatroom)
}

func (s *chatroomService) MigrateChatroomsToUserChatrooms() error {
	return s.repo.MigrateChatroomsToUserChatrooms()
}

func (s *chatroomService) GetChatroomByID(id string) (*domain.Entity_Chatroom, error) {
	return s.repo.GetChatroomByID(id)
}

func (s *chatroomService) GetUserChatrooms(userID string, page, pageSize int) ([]DTO_ChatroomResult, error) {

	rawData, err := s.repo.GetUserChatrooms(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	var chatrooms []DTO_ChatroomResult
	for _, item := range rawData {
		chatroom := DTO_ChatroomResult{
			Title:        item.Chatroom.Title,
			Avatar:       item.Chatroom.Avatar,
			CreatorID:    item.Chatroom.CreatorID,
			CreatedAt:    item.Chatroom.CreatedAt.UTC(),
			UpdatedAt:    item.Chatroom.UpdatedAt.UTC(),
			Type:         item.Chatroom.Type,
			ChatroomID:   item.Chatroom.ID,
			LastReadTime: item.PersonalConfig.LastReadTime,
			UnreadCount:  item.PersonalConfig.UnreadCount,
			IsMuted:      item.PersonalConfig.IsMuted,
			IsHidden:     item.PersonalConfig.IsHidden,
			IsPinned:     item.PersonalConfig.IsPinned,
			BaseID:       item.Chatroom.BaseID,
			BoardID:      item.Chatroom.BoardID,
			UserList:     item.Chatroom.UserList,
			LastMessage: LastMessage{
				Content:      item.LastMessage.Content,
				Timestamp:    item.LastMessage.Timestamp,
				SenderID:     item.LastMessage.SenderID,
				SenderName:   item.LastMessage.SenderName,
				SenderAvatar: item.LastMessage.SenderAvatar,
			},
		}
		chatrooms = append(chatrooms, chatroom)

	}

	return chatrooms, nil
}

func (s *chatroomService) UpdateChatroom(chatroomID string, title string, avatar string) error {
	existingChatroom, err := s.repo.GetChatroomByID(chatroomID)
	if err != nil {
		return errors.New("chatroom not found")
	}
	existingChatroom.Avatar = avatar
	existingChatroom.Title = title

	return s.repo.UpdateChatroom(existingChatroom)
}
func (s *chatroomService) IsUserInChatroom(chatroom *domain.Entity_Chatroom, userID string) bool {
	for _, userList := range chatroom.UserList {
		if userList == userID {
			return true
		}
	}
	return false
}

func (s *chatroomService) SendMessage(dto DTO_SendMessageRequest) (*DTO_MessagesResult, error) {
	chatroomData, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil || chatroomData == nil {
		return nil, errors.New("chatroom does not exist")
	}
	if !s.IsUserInChatroom(chatroomData, dto.SenderID) {
		return nil, errors.New("user is not a member of the chatroom")
	}

	ctime := time.Now()
	utime := ctime
	message := &domain.Entity_Message{
		ChatroomID: dto.ChatroomID,
		SenderID:   dto.SenderID,
		Avatar:     dto.Avatar,
		Nickname:   dto.Nickname,
		Type:       dto.Type,
		Data:       dto.Data,
		CreatedAt:  ctime,
		UpdatedAt:  utime,
	}

	savedMessage, err := s.repo.SendMessage(message)
	if err != nil {
		return nil, err
	}

	chatroomData.LastMessageID = dto.Data

	if err := s.repo.UpdateChatroom(chatroomData); err != nil {
		return nil, err
	}

	messagesResult := DTO_MessagesResult{
		ID:         savedMessage.ID,
		ChatroomID: savedMessage.ChatroomID,
		SenderID:   savedMessage.SenderID,
		Type:       savedMessage.Type,
		Data:       savedMessage.Data,
		CreatedAt:  savedMessage.CreatedAt,
		UpdatedAt:  savedMessage.UpdatedAt,
	}

	return &messagesResult, nil
}

func (s *chatroomService) GetMessages(chatroomID string, page int, pageSize int) ([]*domain.Entity_Message, error) {
	return s.repo.GetMessages(chatroomID, page, pageSize)
}

func (s *chatroomService) GetChatroomUserList(chatroomID string) ([]string, error) {
	return s.repo.GetChatroomUserList(chatroomID)
}

func (s *chatroomService) RemoveUserFromChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	check_chatroom, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil {
		return errors.New("failed to check chatroom ID: ")
		// + err.Error())
	}
	if check_chatroom == nil {
		return errors.New("chatroom does not exist")
	}

	if check_chatroom.CreatorID != dto.UserID {
		return errors.New("no permission")
	}

	_chatroom := domain.SpecifyChatRoomUser{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
	}

	err = s.repo.RemoveUserFromChatroom(_chatroom)
	if err != nil {
		return errors.New("failed to remove user from chatroom: " + err.Error())
	}
	return s.repo.RemoveUserFromChatroom(_chatroom)
}

func (s *chatroomService) LeaveChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	chatroom, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil {
		return errors.New("failed to check chatroom ID: ")
		// + err.Error())
	}
	if chatroom == nil {
		return errors.New("chatroom does not exist")
	}
	_specifyChatRoomUser := domain.SpecifyChatRoomUser{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
	}

	err = s.repo.RemoveUserFromChatroom(_specifyChatRoomUser)
	if err != nil {
		return errors.New("failed to remove user from chatroom: " + err.Error())
	}
	return s.repo.RemoveUserFromChatroom(_specifyChatRoomUser)
}

func (s *chatroomService) AddUserToChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	newConfig := domain.ChatroomPersonalConfig{
		ChatroomID:  dto.ChatroomID,
		UnreadCount: 0,
		IsMuted:     false,
		IsHidden:    false,
		IsPinned:    false,
	}

	_specifyChatRoomUser := domain.SpecifyChatRoomUser{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
	}

	err := s.repo.AddUserToChatroom(_specifyChatRoomUser, &newConfig)
	if err != nil {
		return fmt.Errorf("failed to add user to chatroom: %w", err)
	}

	return nil
}
