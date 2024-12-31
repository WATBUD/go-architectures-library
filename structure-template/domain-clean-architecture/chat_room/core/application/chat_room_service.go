package application

import (
	"chat_room_mod/core/domain/chatroom"
	"errors"
	"time"
)

type chatroomService struct {
	repo chatroom.ChatroomRepository
}

func NewChatroomService(repo chatroom.ChatroomRepository) *chatroomService {
	return &chatroomService{repo: repo}
}

func (s *chatroomService) SetUserChatRoomConfig(dto chatroom.DTO_SetUserChatRoomConfig) error {
	return s.repo.SetUserChatRoomConfig(dto)
}

// func (s *chatroomService) UserMutedChatRoom(dto chatroom.DTO_SetUserChatRoomConfig) error {
// 	return s.repo.SetUserChatRoomConfig(dto)
// }

func (s *chatroomService) CreateChatroom(dto chatroom.DTO_CreateChatroomRequest) (*chatroom.Entity_Chatroom, error) {
	ctime := time.Now()
	utime := ctime
	personalConfig := chatroom.ChatroomPersonalConfig{
		NotificationOn: true,
		UserID:         dto.CreatorID,
		IsMuted:        false,
		IsHidden:       false,
		IsPinned:       false,
	}

	chatroom := &chatroom.Entity_Chatroom{
		Title:         dto.Title,
		CreatorID:     dto.CreatorID,
		BaseID:        dto.BaseID,
		BoardID:       dto.BoardID,
		Avatar:        dto.Avatar,
		LastMsg:       dto.Nikename + "_已創建聊天室",
		LastMsgSender: dto.Nikename,
		LastMsgTime:   utime,
		CreatedAt:     ctime,
		UpdatedAt:     utime,
		Type:          "private",
		PersonalConfig: []*chatroom.ChatroomPersonalConfig{
			&personalConfig,
		},
	}
	return s.repo.CreateChatroom(chatroom)
}

func (s *chatroomService) GetChatroomByID(id string) (*chatroom.Entity_Chatroom, error) {
	return s.repo.GetChatroomByID(id)
}

func (s *chatroomService) GetChatrooms(userID string) ([]*chatroom.Entity_Chatroom, error) {
	return s.repo.GetChatrooms(userID)
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
func (s *chatroomService) IsUserInChatroom(chatroom *chatroom.Entity_Chatroom, userID string) bool {
	for _, config := range chatroom.PersonalConfig {
		if config.UserID == userID {
			return true
		}
	}
	return false
}

func (s *chatroomService) SendMessage(dto chatroom.DTO_SendMessageRequest) (*chatroom.Entity_Message, error) {
	chatroomData, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil || chatroomData == nil {
		return nil, errors.New("chatroom does not exist")
	}
	if !s.IsUserInChatroom(chatroomData, dto.SenderID) {
		return nil, errors.New("user is not a member of the chatroom")
	}

	ctime := time.Now()
	utime := ctime
	message := &chatroom.Entity_Message{
		ChatroomID: dto.ChatroomID,
		SenderID:   dto.SenderID,
		Type:       dto.Type,
		Data:       dto.Data,
		CreatedAt:  ctime,
		UpdatedAt:  utime,
	}

	savedMessage, err := s.repo.SendMessage(message)
	if err != nil {
		return nil, err
	}

	chatroomData.LastMsg = dto.Data
	chatroomData.LastMsgTime = ctime
	chatroomData.LastMsgSender = dto.SenderID

	if err := s.repo.UpdateChatroom(chatroomData); err != nil {
		return nil, err
	}

	return savedMessage, nil
}

func (s *chatroomService) GetMessages(chatroomID string) ([]*chatroom.Entity_Message, error) {
	return s.repo.GetMessages(chatroomID)
}

func (s *chatroomService) GetChatroomPersonalConfig(chatroomID string) ([]*chatroom.ChatroomPersonalConfig, error) {
	return s.repo.GetChatroomPersonalConfig(chatroomID)
}

func (s *chatroomService) RemoveUserFromChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error {
	chatroom, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil {
		return errors.New("failed to check chatroom ID: ")
		// + err.Error())
	}
	if chatroom == nil {
		return errors.New("chatroom does not exist")
	}

	if chatroom.CreatorID != dto.UserID {
		return errors.New("no permission")
	}

	err = s.repo.RemoveUserFromChatroom(dto)
	if err != nil {
		return errors.New("failed to remove user from chatroom: " + err.Error())
	}
	return s.repo.RemoveUserFromChatroom(dto)
}

func (s *chatroomService) AddUserToChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error {
	chatroomData, err := s.repo.GetChatroomByID(dto.ChatroomID)
	if err != nil {
		return errors.New("failed to check chatroom ID: ")
		// + err.Error())
	}
	if chatroomData == nil {
		return errors.New("chatroom does not exist")
	}
	for _, config := range chatroomData.PersonalConfig {
		if config.UserID == dto.UserID {
			return errors.New("user already exists in chatroom")
		}
	}
	personalConfig := append(chatroomData.PersonalConfig, &chatroom.ChatroomPersonalConfig{
		UserID:         dto.UserID,
		UnreadCount:    0,
		NotificationOn: true,
		IsMuted:        false,
		IsHidden:       false,
		IsPinned:       false,
	})

	err = s.repo.AddUserToChatroom(dto.ChatroomID, personalConfig)
	if err != nil {
		return errors.New("failed to add user to chatroom: " + err.Error())
	}
	return s.repo.AddUserToChatroom(dto.ChatroomID, personalConfig)
}
