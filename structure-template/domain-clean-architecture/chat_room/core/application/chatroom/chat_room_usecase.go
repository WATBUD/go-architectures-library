package application_chatroom

import (
	"context"
	"fmt"
	"time"

	"mai.today/user"
)

type ChatroomUsecase interface {
	CreateChatroom(dto DTO_CreateChatroomRequest) (*DTO_ChatroomResult, error)
	GetChatroomByID(id string) (*DTO_ChatroomResult, error)
	UpdateChatroom(chatroomID string, Title string, Avatar string) error
	SendMessage(dto DTO_SendMessageRequest) (*DTO_MessagesResult, error)
	GetMessages(chatroomID string, page int, pageSize int) ([]*DTO_MessagesResult, error)
	RemoveUserFromChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	LeaveChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	AddUserToChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error
	UserMutedChatRoom(dto DTO_MuteUserRequest) error
	UserPinnedChatRoom(dto DTO_PinUserRequest) error
	UserHideChatRoom(dto DTO_HiddenUserRequest) error

	GetUserChatrooms(userID string, page int, pageSize int) ([]*DTO_ChatroomResult, error)
	GetChatroomMembers(chatroomID string) ([]*DTO_MembersResult, error)
}

type chatroomUsecase struct {
	service ChatroomService
}

func NewChatroomUsecase(service ChatroomService) ChatroomUsecase {
	return &chatroomUsecase{service: service}
}

func (u *chatroomUsecase) UserPinnedChatRoom(dto DTO_PinUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID

	dtoConfig := DTO_SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  "is_pinned",
		ParamValue: dto.Pin,
	}

	return u.service.SetUserChatRoomConfig(dtoConfig)
}

func (u *chatroomUsecase) UserHideChatRoom(dto DTO_HiddenUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID
	dtoConfig := DTO_SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  "is_hidden",
		ParamValue: dto.IsHidden,
	}

	return u.service.SetUserChatRoomConfig(dtoConfig)
}

func (u *chatroomUsecase) UserMutedChatRoom(dto DTO_MuteUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID

	dtoConfig := DTO_SetUserChatRoomConfig{
		ChatroomID: dto.ChatroomID,
		UserID:     dto.UserID,
		ParamName:  "is_muted",
		ParamValue: dto.Mute,
	}

	return u.service.SetUserChatRoomConfig(dtoConfig)
}

func (u *chatroomUsecase) CreateChatroom(dto DTO_CreateChatroomRequest) (*DTO_ChatroomResult, error) {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.CreatorID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return nil, fmt.Errorf("UID %s not found", dto.CreatorID)
	}
	dto.CreatorID = userInfo.UserInfo.FirebaseUID
	dto.Nickname = userInfo.UserInfo.Nickname
	dto.Avatar = userInfo.UserInfo.Avatar.Path

	createdChatroom, err := u.service.CreateChatroom(dto)
	if err != nil {
		return nil, err
	}

	return &DTO_ChatroomResult{
		ChatroomID: createdChatroom.ID,
		BaseID:     createdChatroom.BaseID,
		BoardID:    createdChatroom.BoardID,
		Title:      createdChatroom.Title,
		LastMessage: LastMessage{
			Content:      "",
			Timestamp:    time.Time{},
			SenderID:     dto.CreatorID,
			SenderName:   dto.Nickname,
			SenderAvatar: dto.Avatar,
		},
		Avatar:    createdChatroom.Avatar,
		CreatorID: createdChatroom.CreatorID,
		CreatedAt: createdChatroom.CreatedAt,
		UpdatedAt: createdChatroom.UpdatedAt,
		UserList:  createdChatroom.UserList,
		Type:      createdChatroom.Type,
	}, nil
}

func (u *chatroomUsecase) GetChatroomByID(id string) (*DTO_ChatroomResult, error) {

	dto, err := u.service.GetChatroomByID(id)
	if err != nil {
		return nil, err
	}

	return &DTO_ChatroomResult{
		ChatroomID: dto.ID,
		BaseID:     dto.BaseID,
		BoardID:    dto.BoardID,
		Title:      dto.Title,
		LastMessage: LastMessage{
			Content:   "",
			Timestamp: time.Time{},
			SenderID:  dto.CreatorID,
			//SenderName:   dto.Nickname,
			SenderAvatar: dto.Avatar,
		},
		// LastMsg:       dto.LastMsg,
		// LastMsgTime:   dto.LastMsgTime,
		// LastMsgSender: dto.LastMsgSender,
		Avatar:    dto.Avatar,
		CreatorID: dto.CreatorID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		UserList:  dto.UserList,
		Type:      dto.Type,
	}, nil
}

func (u *chatroomUsecase) GetUserChatrooms(userID string, page int, pageSize int) ([]*DTO_ChatroomResult, error) {
	aggregation, err := u.service.GetUserChatrooms(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), userID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return nil, fmt.Errorf("UID %s not found", userID)
	}
	// dto.CreatorID = userInfo.UserInfo.FirebaseUID
	// dto.Nickname = userInfo.UserInfo.Nickname
	// dto.Avatar = userInfo.UserInfo.Avatar.Path

	data := make([]*DTO_ChatroomResult, len(aggregation))
	for i, item := range aggregation {
		// userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), item.LastMessage.SenderID)
		// if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		// 	return nil, fmt.Errorf("UID %s not found", item.LastMessage.SenderID)
		// }
		data[i] = &DTO_ChatroomResult{
			BaseID:  item.BaseID,
			BoardID: item.BaseID,
			Title:   item.Title,
			LastMessage: LastMessage{
				Content:      item.LastMessage.Content,
				Timestamp:    item.LastMessage.Timestamp,
				SenderID:     item.LastMessage.SenderID,
				SenderName:   item.LastMessage.SenderName,
				SenderAvatar: item.LastMessage.SenderAvatar,
			},
			Avatar:       item.Avatar,
			CreatorID:    item.CreatorID,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
			UserList:     item.UserList,
			Type:         item.Type,
			ChatroomID:   item.ChatroomID,
			LastReadTime: item.LastReadTime,
			UnreadCount:  item.UnreadCount,
			IsMuted:      item.IsMuted,
			IsHidden:     item.IsHidden,
			IsPinned:     item.IsPinned,
		}
	}

	return data, nil
}

func (u *chatroomUsecase) GetMessages(chatroomID string, page int, pageSize int) ([]*DTO_MessagesResult, error) {
	_messages, err := u.service.GetMessages(chatroomID, page, pageSize)
	if err != nil {
		return nil, err
	}

	var resultMessages []*DTO_MessagesResult

	// Check if _messages is empty
	if len(_messages) == 0 {
		return []*DTO_MessagesResult{}, nil
	}

	for _, message := range _messages {
		resultMessage := &DTO_MessagesResult{
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

func (u *chatroomUsecase) GetChatroomMembers(chatroomID string) ([]*DTO_MembersResult, error) {

	userList, err := u.service.GetChatroomUserList(chatroomID)
	if err != nil {
		return nil, err
	}
	var membersResult []*DTO_MembersResult
	for _, _userList := range userList {
		userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), _userList)

		resultMessage := &DTO_MembersResult{
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

func (u *chatroomUsecase) SendMessage(dto DTO_SendMessageRequest) (*DTO_MessagesResult, error) {

	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.SenderID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return nil, fmt.Errorf("UID %s not found", dto.SenderID)
	}

	entityMessage := &DTO_SendMessageRequest{
		ChatroomID: dto.ChatroomID,
		SenderID:   userInfo.UserInfo.FirebaseUID,
		Avatar:     userInfo.UserInfo.Avatar.Path,
		Nickname:   userInfo.UserInfo.Nickname,
		Type:       dto.Type,
		Data:       dto.Data,
	}
	createdMessage, err := u.service.SendMessage(*entityMessage)
	if err != nil {
		return nil, err
	}

	dtoMessage := &DTO_MessagesResult{
		ID:         createdMessage.ID,
		ChatroomID: createdMessage.ChatroomID,
		SenderID:   createdMessage.SenderID,
		Avatar:     createdMessage.Avatar,
		Nickname:   createdMessage.Nickname,
		Type:       createdMessage.Type,
		Data:       createdMessage.Data,
		CreatedAt:  createdMessage.CreatedAt,
		UpdatedAt:  createdMessage.UpdatedAt,
	}

	return dtoMessage, nil
}

func (u *chatroomUsecase) RemoveUserFromChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	return u.service.RemoveUserFromChatroom(dto)
}
func (u *chatroomUsecase) LeaveChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	return u.service.LeaveChatroom(dto)
}

func (u *chatroomUsecase) AddUserToChatroom(dto DTO_AddOrRemoveChatRoomUserRequest) error {
	userInfo := user.Instance().FindUserByFireBaseUID(context.Background(), dto.UserID)
	if userInfo == nil || userInfo.UserInfo.FirebaseUID == "" {
		return fmt.Errorf("UID %s not found", dto.UserID)
	}
	dto.UserID = userInfo.UserInfo.FirebaseUID
	return u.service.AddUserToChatroom(dto)
}
