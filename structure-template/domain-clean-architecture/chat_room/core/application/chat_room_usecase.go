package application

import (
	"chat_room_mod/core/domain/chatroom"
	"context"
	"fmt"

	"mai.today/user"
)

type ChatroomUsecase interface {
	CreateChatroom(dto chatroom.DTO_CreateChatroomRequest) (*chatroom.Entity_Chatroom, error)
	GetChatroomByID(id string) (*chatroom.Entity_Chatroom, error)
	GetChatrooms(userID string) ([]*chatroom.Entity_Chatroom, error)
	UpdateChatroom(chatroomID string, Title string, Avatar string) error
	SendMessage(dto chatroom.DTO_SendMessageRequest) (*chatroom.Entity_Message, error)
	GetMessages(chatroomID string) ([]*chatroom.DTO_MessagesResult, error)
	RemoveUserFromChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error
	AddUserToChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error
	UserMutedChatRoom(dto chatroom.DTO_MuteUserRequest) error
	UserPinnedChatRoom(dto chatroom.DTO_PinUserRequest) error
	GetChatroomMembers(chatroomID string) ([]*chatroom.DTO_MembersResult, error)
}

type chatroomUsecase struct {
	service *chatroomService
}

func NewChatroomUsecase(service *chatroomService) ChatroomUsecase {
	return &chatroomUsecase{service: service}
}

func (u *chatroomUsecase) UserPinnedChatRoom(dto chatroom.DTO_PinUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID

	dtoConfig := chatroom.DTO_SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  "is_pinned",
		ParamValue: dto.Pin,
	}

	return u.service.SetUserChatRoomConfig(dtoConfig)
}

func (u *chatroomUsecase) UserMutedChatRoom(dto chatroom.DTO_MuteUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID

	dtoConfig := chatroom.DTO_SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  "is_muted",
		ParamValue: dto.Mute,
	}

	return u.service.SetUserChatRoomConfig(dtoConfig)
}

func (u *chatroomUsecase) CreateChatroom(dto chatroom.DTO_CreateChatroomRequest) (*chatroom.Entity_Chatroom, error) {

	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.CreatorID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return nil, fmt.Errorf("UID %s not found", dto.CreatorID)
	}
	dto.CreatorID = userInfo.UserInfo.FirebaseUID
	dto.Nikename = userInfo.UserInfo.Nickname

	return u.service.CreateChatroom(dto)
}

func (u *chatroomUsecase) GetChatroomByID(id string) (*chatroom.Entity_Chatroom, error) {
	return u.service.GetChatroomByID(id)
}

func (u *chatroomUsecase) GetChatrooms(userID string) ([]*chatroom.Entity_Chatroom, error) {

	return u.service.GetChatrooms(userID)
}

func (u *chatroomUsecase) GetMessages(chatroomID string) ([]*chatroom.DTO_MessagesResult, error) {
	_messages, err := u.service.GetMessages(chatroomID)
	if err != nil {
		return nil, err
	}
	var resultMessages []*chatroom.DTO_MessagesResult

	// Check if _messages is empty
	if len(_messages) == 0 {
		return []*chatroom.DTO_MessagesResult{}, nil
	}
	for _, message := range _messages {
		resultMessage := &chatroom.DTO_MessagesResult{
			ID:         message.ID,
			ChatroomID: message.ChatroomID,
			SenderID:   message.SenderID,
			Nickname:   "不存在的使用者",
			Type:       message.Type,
			Data:       message.Data,
			CreatedAt:  message.CreatedAt,
			UpdatedAt:  message.UpdatedAt,
		}
		userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), message.SenderID)

		if userInfo != nil {
			resultMessage.Nickname = userInfo.UserInfo.Nickname
		}

		resultMessages = append(resultMessages, resultMessage)
	}

	return resultMessages, nil
}

func (u *chatroomUsecase) GetChatroomMembers(chatroomID string) ([]*chatroom.DTO_MembersResult, error) {
	personalConfigs, err := u.service.GetChatroomPersonalConfig(chatroomID)
	if err != nil {
		return nil, err
	}
	var membersResult []*chatroom.DTO_MembersResult
	for _, personalConfig := range personalConfigs {
		userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), personalConfig.UserID)

		resultMessage := &chatroom.DTO_MembersResult{
			Nickname: "不存在的使用者",
		}
		if userInfo != nil {
			resultMessage.Nickname = userInfo.UserInfo.Nickname
			resultMessage.Avatar = userInfo.UserInfo.Avatar.Path
			resultMessage.UserID = userInfo.UserInfo.FirebaseUID
		}

		membersResult = append(membersResult, resultMessage)
	}

	return membersResult, nil
}

func (u *chatroomUsecase) UpdateChatroom(chatroomID string, Title string, Avatar string) error {
	return u.service.UpdateChatroom(chatroomID, Title, Avatar)
}

func (u *chatroomUsecase) SendMessage(dto chatroom.DTO_SendMessageRequest) (*chatroom.Entity_Message, error) {

	return u.service.SendMessage(dto)
}

func (u *chatroomUsecase) RemoveUserFromChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error {
	return u.service.RemoveUserFromChatroom(dto)
}

func (u *chatroomUsecase) AddUserToChatroom(dto chatroom.DTO_AddOrRemoveChatRoomUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID
	return u.service.AddUserToChatroom(dto)
}
