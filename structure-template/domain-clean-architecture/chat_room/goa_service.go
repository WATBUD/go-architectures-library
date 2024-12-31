package chat_room_mod

import (
	"chat_room_mod/adapters/repository"
	application "chat_room_mod/core/application/chatroom"
	domain "chat_room_mod/core/domain/chatroom"
	"context"
	"errors"
	"sync"
	"time"

	goa_chat_room "mai.today/api/gen/chatroom"
	"mai.today/api/gen/http/chatroom/server"
	"mai.today/authentication"
	"mai.today/database/mongodb"
	usercontext "mai.today/foundation/context/user"
	"mai.today/realtime"
)

func NewChatroomServices(repo domain.ChatroomRepository) *ChatroomService {
	return &ChatroomService{
		applicationUsecase: application.NewChatroomUsecase(application.NewChatroomService(repo)),
	}
}

func (cs *ChatroomService) CreateChatroom(ctx context.Context, p *goa_chat_room.CreateChatroomPayload) (*goa_chat_room.CreateChatroomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	var boardID string
	if p.BoardID != nil {
		boardID = *p.BoardID
	}

	var avatar string
	if p.Avatar != nil {
		avatar = *p.Avatar
	}

	dto := application.DTO_CreateChatroomRequest{
		Title:     p.Title,
		BaseID:    p.BaseID,
		BoardID:   boardID,
		Type:      "private",
		Avatar:    avatar,
		CreatorID: uid,
	}

	createChatroom, err := cs.applicationUsecase.CreateChatroom(dto)
	if err != nil {
		return nil, err
	}
	// err = cs.realtime.SubscribeUser(ctx, createChatroom.ID, uid)
	// if err != nil {
	// 	return nil, err
	// }
	resultData := &goa_chat_room.CreateChatroomResult{
		Command:   &goa_chat_room.Command{Type: "CreateChatroom"},
		Timestamp: createChatroom.CreatedAt.Unix(),
		Data: &goa_chat_room.Chatroom{
			ChatroomID: createChatroom.ChatroomID,
			BaseID:     createChatroom.BaseID,
			BoardID:    createChatroom.BoardID,
			Title:      createChatroom.Title,
			Avatar:     createChatroom.Avatar,
			CreatedAt:  createChatroom.CreatedAt.UTC(),
			UpdatedAt:  createChatroom.UpdatedAt.UTC(),
			Type:       createChatroom.Type,
			//Members:    createChatroom.Members,
		},
	}
	_, err = cs.realtime.PublishToUser(ctx, uid, server.NewCreateChatroomOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}
	return resultData, nil
}

func (cs *ChatroomService) RemoveUserFromChatroom(ctx context.Context, p *goa_chat_room.RemoveUserFromChatroomPayload) (*goa_chat_room.RemoveUserFromChatroomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_AddOrRemoveChatRoomUserRequest{
		ChatroomID: p.ChatroomID,
		UserID:     uid,
	}

	err := cs.applicationUsecase.RemoveUserFromChatroom(dto)
	if err != nil {
		return nil, err
	}

	return &goa_chat_room.RemoveUserFromChatroomResult{
		Command:   &goa_chat_room.Command{Type: "RemoveUserFromChatroom"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.RemoveUser{
			ChatroomID: p.ChatroomID,
			Message:    "User removed successfully",
		},
	}, nil
}

func (cs *ChatroomService) UpdateChatroom(ctx context.Context, p *goa_chat_room.UpdateChatroomPayload) (*goa_chat_room.UpdateChatroomResult, error) {
	chatroomID := p.ChatroomID
	title := p.Title
	avatar := p.Avatar

	err := cs.applicationUsecase.UpdateChatroom(chatroomID, title, avatar)
	if err != nil {
		if err.Error() == "chatroom not found" {
			return nil, errors.New("chatroom not found")
		}
		return nil, errors.New("internal server error")
	}

	updatedChatroom, err := cs.applicationUsecase.GetChatroomByID(chatroomID)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	resultData := &goa_chat_room.UpdateChatroomResult{
		Command:   &goa_chat_room.Command{Type: "UpdateChatroom"},
		Timestamp: updatedChatroom.UpdatedAt.Unix(),
		Data: &goa_chat_room.Chatroom{
			ChatroomID: updatedChatroom.ChatroomID,
			BaseID:     updatedChatroom.BaseID,
			Title:      updatedChatroom.Title,
			Avatar:     updatedChatroom.Avatar,
			CreatedAt:  updatedChatroom.CreatedAt.UTC(),
			UpdatedAt:  updatedChatroom.UpdatedAt.UTC(),
			BoardID:    updatedChatroom.BoardID,
		},
	}

	members, err := cs.applicationUsecase.GetChatroomMembers(p.ChatroomID)
	if err != nil {
		return nil, err
	}

	var userIDs []string

	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	_, err = cs.realtime.BroadcastToUsers(ctx, userIDs, server.NewUpdateChatroomOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}

	return resultData, nil
}

func (cs *ChatroomService) GetMessages(ctx context.Context, p *goa_chat_room.GetMessagesPayload) (*goa_chat_room.GetMessagesResult, error) {
	//uid, _ := usercontext.GetUserID(ctx)

	messages, err := cs.applicationUsecase.GetMessages(p.ChatroomID, p.Page, p.PageSize)
	if err != nil {
		return nil, err
	}

	var messageResults []*goa_chat_room.MessageResult
	for _, msg := range messages {
		messageResults = append(messageResults, &goa_chat_room.MessageResult{
			Command:   &goa_chat_room.Command{Type: "GetMessages"},
			Timestamp: time.Now().Unix(),
			Data: &goa_chat_room.Message{
				ID:         msg.ID,
				ChatroomID: msg.ChatroomID,
				SenderID:   msg.SenderID,
				Type:       msg.Type,
				Data:       msg.Data,
				Avatar:     msg.Avatar,
				Nickname:   msg.Nickname,
				CreatedAt:  msg.CreatedAt.UTC(),
				UpdatedAt:  msg.UpdatedAt.UTC(),
			},
		})
	}
	messageData := &goa_chat_room.Messages{
		Messages: messageResults,
	}

	resultData := &goa_chat_room.GetMessagesResult{
		Command:   &goa_chat_room.Command{Type: "GetMessages"},
		Timestamp: time.Now().Unix(),
		Data:      messageData,
	}
	return resultData, nil
}

func (cs *ChatroomService) ListChatrooms(ctx context.Context, p *goa_chat_room.ListChatroomsPayload) (*goa_chat_room.ListChatroomsResult, error) {

	uid, _ := usercontext.GetUserID(ctx)
	chatrooms, err := cs.applicationUsecase.GetUserChatrooms(uid, p.Page, p.PageSize)
	if err != nil {
		return nil, err
	}

	var chatroomCollection []*goa_chat_room.Chatroom
	for _, chatItem := range chatrooms {
		chatroomCollection = append(chatroomCollection, &goa_chat_room.Chatroom{
			ChatroomID: chatItem.ChatroomID,
			BaseID:     chatItem.BaseID,
			BoardID:    chatItem.BoardID,
			Title:      chatItem.Title,
			Avatar:     chatItem.Avatar,
			CreatedAt:  chatItem.CreatedAt.UTC(),
			UpdatedAt:  chatItem.UpdatedAt.UTC(),
			Type:       chatItem.Type,
			IsMuted:    &chatItem.IsMuted,
			IsHidden:   &chatItem.IsHidden,
			IsPinned:   &chatItem.IsPinned,
		})
	}

	result := &goa_chat_room.ListChatroomsResult{
		Data: chatroomCollection,
	}
	return result, nil
}

func (cs *ChatroomService) SendMessage(ctx context.Context, p *goa_chat_room.SendMessagePayload) (*goa_chat_room.SendMessageResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_SendMessageRequest{
		ChatroomID: p.ChatroomID,
		SenderID:   uid,
		Type:       p.Type,
		Data:       p.Data,
	}

	message, err := cs.applicationUsecase.SendMessage(dto)
	if err != nil {
		return nil, err
	}
	// messageJSON, err := json.Marshal(message)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to marshal message: %w", err)
	// }
	resultData := &goa_chat_room.SendMessageResult{
		Command:   &goa_chat_room.Command{Type: "SendMessage"},
		Timestamp: message.UpdatedAt.Unix(),
		Data: &goa_chat_room.Message{
			ID:         message.ID,
			ChatroomID: message.ChatroomID,
			SenderID:   message.SenderID,
			Type:       message.Type,
			Data:       message.Data,
			CreatedAt:  message.CreatedAt.UTC(),
			UpdatedAt:  message.UpdatedAt.UTC(),
		},
	}

	members, err := cs.applicationUsecase.GetChatroomMembers(p.ChatroomID)
	if err != nil {
		return nil, err
	}

	var userIDs []string

	for _, member := range members {
		userIDs = append(userIDs, member.UserID)
	}
	_, err = cs.realtime.BroadcastToUsers(ctx, userIDs, server.NewSendMessageOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}

	return resultData, nil
}

func (cs *ChatroomService) GetChatroomMembers(ctx context.Context, p *goa_chat_room.GetChatroomMembersPayload) (*goa_chat_room.GetChatroomMembersResult, error) {

	members, err := cs.applicationUsecase.GetChatroomMembers(p.ChatroomID)
	if err != nil {
		return nil, err
	}

	var memberResults []*goa_chat_room.ChatroomMemberResult
	for _, member := range members {
		memberResults = append(memberResults, &goa_chat_room.ChatroomMemberResult{
			UserID:   member.UserID,
			Avatar:   member.Avatar,
			Nickname: member.Nickname,
		})
	}

	resultData := &goa_chat_room.GetChatroomMembersResult{
		Command:   &goa_chat_room.Command{Type: "GetChatroomMembers"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.ChatroomMembers{
			Members: memberResults,
		},
	}

	uid, _ := usercontext.GetUserID(ctx)

	_, err = cs.realtime.PublishToUser(ctx, uid, server.NewGetChatroomMembersOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}

	return resultData, nil
}

func (cs *ChatroomService) UserPinnedChatRoom(ctx context.Context, p *goa_chat_room.UserPinnedChatRoomPayload) (*goa_chat_room.UserPinnedChatRoomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_PinUserRequest{
		ChatroomID: p.ChatroomID,
		UserID:     uid,
		Pin:        &p.Pin,
	}

	err := cs.applicationUsecase.UserPinnedChatRoom(dto)
	if err != nil {
		return nil, err
	}

	resultData := &goa_chat_room.UserPinnedChatRoomResult{
		Command:   &goa_chat_room.Command{Type: "UserPinnedChatRoom"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.UserPinnedChatroom{
			ChatroomID: p.ChatroomID,
			Pin:        p.Pin,
		},
	}

	_, err = cs.realtime.PublishToUser(ctx, uid, server.NewUserPinnedChatRoomOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}
	return resultData, nil
}

func (cs *ChatroomService) UserHideChatRoom(ctx context.Context, p *goa_chat_room.UserHideChatRoomPayload) (*goa_chat_room.UserHideChatRoomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_PinUserRequest{
		ChatroomID: p.ChatroomID,
		UserID:     uid,
		Pin:        &p.Hide,
	}
	err := cs.applicationUsecase.UserPinnedChatRoom(dto)
	if err != nil {
		return nil, err
	}

	resultData := &goa_chat_room.UserHideChatRoomResult{
		Command:   &goa_chat_room.Command{Type: "UserPinnedChatRoom"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.UserHideChatroom{
			ChatroomID: p.ChatroomID,
			UserID:     uid,
			Hide:       p.Hide,
		},
	}

	_, err = cs.realtime.PublishToUser(ctx, uid, server.NewUserHideChatRoomOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}
	return resultData, nil
}

func (cs *ChatroomService) ReceiveUserHideChatRoom(context.Context, *goa_chat_room.ReceiveUserHideChatRoomPayload) (res *goa_chat_room.UserHideChatRoomResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) UserMutedChatRoom(ctx context.Context, p *goa_chat_room.UserMutedChatRoomPayload) (*goa_chat_room.UserMutedChatRoomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_MuteUserRequest{
		ChatroomID: p.ChatroomID,
		UserID:     uid,
		Mute:       &p.Mute,
	}

	err := cs.applicationUsecase.UserMutedChatRoom(dto)
	if err != nil {
		return nil, err
	}

	resultData := &goa_chat_room.UserMutedChatRoomResult{
		Command:   &goa_chat_room.Command{Type: "UserMutedChatRoom"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.UserMutedChatroom{
			ChatroomID: p.ChatroomID,
			UserID:     uid,
			Mute:       p.Mute,
		},
	}
	_, err = cs.realtime.PublishToUser(ctx, uid, server.NewUserMutedChatRoomOKResponseBody(resultData))
	if err != nil {
		return nil, err
	}
	return resultData, nil
}

func (cs *ChatroomService) AddUserToChatroom(ctx context.Context, p *goa_chat_room.AddUserToChatroomPayload) (*goa_chat_room.AddUserToChatroomResult, error) {
	uid, _ := usercontext.GetUserID(ctx)
	dto := application.DTO_AddOrRemoveChatRoomUserRequest{
		ChatroomID: p.ChatroomID,
		UserID:     uid,
	}

	err := cs.applicationUsecase.AddUserToChatroom(dto)
	if err != nil {
		return nil, err
	}

	return &goa_chat_room.AddUserToChatroomResult{
		Command:   &goa_chat_room.Command{Type: "AddUserToChatroom"},
		Timestamp: time.Now().Unix(),
		Data: &goa_chat_room.AddUser{
			ChatroomID: p.ChatroomID,
			Message:    "User added successfully",
		},
	}, nil
}

func (cs *ChatroomService) ReceiveAddUserToChatroom(context.Context, *goa_chat_room.ReceiveAddUserToChatroomPayload) (res *goa_chat_room.AddUserToChatroomResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) ReceiveGetMessages(context.Context, *goa_chat_room.ReceiveGetMessagesPayload) (res *goa_chat_room.GetMessagesResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) ReceiveSendMessage(context.Context, *goa_chat_room.ReceiveSendMessagePayload) (res *goa_chat_room.SendMessageResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) ReceiveCreateChatroom(context.Context, *goa_chat_room.ReceiveCreateChatroomPayload) (res *goa_chat_room.CreateChatroomResult, err error) {
	panic("not implemented")
}
func (cs *ChatroomService) ReceiveRemoveUserFromChatroom(context.Context, *goa_chat_room.ReceiveRemoveUserFromChatroomPayload) (res *goa_chat_room.RemoveUserFromChatroomResult, err error) {
	panic("not implemented")
}
func (cs *ChatroomService) ReceiveUpdateChatroom(context.Context, *goa_chat_room.ReceiveUpdateChatroomPayload) (res *goa_chat_room.UpdateChatroomResult, err error) {
	panic("not implemented")
}
func (cs *ChatroomService) ReceiveListChatrooms(context.Context, *goa_chat_room.ReceiveListChatroomsPayload) (res *goa_chat_room.ListChatroomsResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) ReceiveUserMutedChatRoom(context.Context, *goa_chat_room.ReceiveUserMutedChatRoomPayload) (res *goa_chat_room.UserMutedChatRoomResult, err error) {
	panic("not implemented")
}
func (cs *ChatroomService) ReceiveUserPinnedChatRoom(context.Context, *goa_chat_room.ReceiveUserPinnedChatRoomPayload) (res *goa_chat_room.UserPinnedChatRoomResult, err error) {
	panic("not implemented")
}

func (cs *ChatroomService) ReceiveGetChatroomMembers(context.Context, *goa_chat_room.ReceiveGetChatroomMembersPayload) (res *goa_chat_room.GetChatroomMembersResult, err error) {
	panic("not implemented")
}

var (
	once     sync.Once
	instance *ChatroomService
)

func Instance() *ChatroomService {
	once.Do(func() {
		instance = newChatroomService()
	})
	return instance
}

type ChatroomService struct {
	realtime realtime.Realtime
	goa_chat_room.Auther
	applicationUsecase application.ChatroomUsecase
	applicationService application.ChatroomService
}

func newChatroomService() *ChatroomService {
	applicationService := application.NewChatroomService(repository.NewMongoChatroomRepository(mongodb.Instance()))
	applicationUsecase := application.NewChatroomUsecase(applicationService)
	return &ChatroomService{
		realtime:           realtime.Instance(),
		Auther:             authentication.Instance(),
		applicationService: applicationService,
		applicationUsecase: applicationUsecase,
	}
}
