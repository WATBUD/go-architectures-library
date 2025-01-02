package repository

import (
	domain "chat_room_mod/core/domain/chatroom"
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoChatroomRepository struct {
	chatroomCollection     *mongo.Collection
	messageCollection      *mongo.Collection
	userChatroomCollection *mongo.Collection
	client                 *mongo.Client
}

func NewMongoChatroomRepository(client *mongo.Client) domain.ChatroomRepository {
	chatroomCollection := client.Database("mai_dev").Collection("chatrooms")
	messageCollection := client.Database("mai_dev").Collection("messages")
	userChatroomCollection := client.Database("mai_dev").Collection("user_chatrooms")

	return &MongoChatroomRepository{
		client:                 client,
		chatroomCollection:     chatroomCollection,
		messageCollection:      messageCollection,
		userChatroomCollection: userChatroomCollection,
	}
}

func (r *MongoChatroomRepository) CreateChatroom(_chatroom *domain.Entity_Chatroom) (*domain.Entity_Chatroom, error) {
	newID := primitive.NewObjectID()
	_chatroom.ID = newID.Hex()
	_, err := r.chatroomCollection.InsertOne(context.TODO(), _chatroom)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return nil, err
	}
	personalConfig := domain.ChatroomPersonalConfig{
		ChatroomID:   _chatroom.ID,
		LastReadTime: time.Now(),
		UnreadCount:  0,
		IsMuted:      false,
		IsHidden:     false,
		IsPinned:     false,
	}
	err = r.AddChatroomToUser(_chatroom.CreatorID, personalConfig)
	if err != nil {
		return nil, err
	}
	return _chatroom, nil
}

func (r *MongoChatroomRepository) AddChatroomToUser(userID string, personalConfig domain.ChatroomPersonalConfig) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$push": bson.M{
			"chatroom_list": personalConfig,
		},
	}
	_, err := r.userChatroomCollection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Printf("Error updating user chatroom list: %v\n", err)
		return err
	}
	return nil
}

func (r *MongoChatroomRepository) MigrateChatroomsToUserChatrooms() error {
	cur, err := r.chatroomCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(context.TODO())

	userChatroomMap := make(map[string][]domain.ChatroomPersonalConfig)
	for cur.Next(context.TODO()) {
		var _chatroom struct {
			ID             string                          `bson:"_id"`
			PersonalConfig []domain.ChatroomPersonalConfig `bson:"personal_config"`
		}
		if err := cur.Decode(&_chatroom); err != nil {
			return err
		}

		// for _, config := range _chatroom.PersonalConfig {
		// 	config.ChatroomID = _chatroom.ID
		// 	userChatroomMap[config.UserID] = append(userChatroomMap[config.UserID], config)
		// }
	}

	for userID, chatroomList := range userChatroomMap {
		userChatroom := domain.Entity_UserChatroom{
			UserID:       userID,
			ChatroomList: chatroomList,
		}

		_, err := r.userChatroomCollection.InsertOne(context.TODO(), userChatroom)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MongoChatroomRepository) GetChatroomByID(id string) (*domain.Entity_Chatroom, error) {
	var _chatroom domain.Entity_Chatroom
	filter := bson.M{"_id": id}
	err := r.chatroomCollection.FindOne(context.TODO(), filter).Decode(&_chatroom)
	if err != nil {
		return nil, err
	}
	return &_chatroom, nil
}

func (r *MongoChatroomRepository) UpdateChatroom(chatroom *domain.Entity_Chatroom) error {
	filter := bson.M{"_id": chatroom.ID}
	update := bson.M{"$set": chatroom}
	_, err := r.chatroomCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoChatroomRepository) SendMessage(message *domain.Entity_Message) (*domain.Entity_Message, error) {

	newID := primitive.NewObjectID()
	message.ID = newID.Hex()
	result, err := r.messageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return nil, err
	}
	fmt.Printf("result: %v\n", result)
	return message, nil
}

func (r *MongoChatroomRepository) GetMessages(chatroomID string, page int, pageSize int) ([]*domain.Entity_Message, error) {
	//DeleteMessagesNotInChatrooms(r)

	filter := bson.M{"chatroom_id": chatroomID}
	options := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(pageSize)).
		SetSkip(int64((page - 1) * pageSize))

	cur, err := r.messageCollection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var messages []*domain.Entity_Message
	for cur.Next(context.TODO()) {
		var _message domain.Entity_Message
		err := cur.Decode(&_message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &_message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func DeleteMessagesNotInChatrooms(r *MongoChatroomRepository) error {
	// Find all chatroom IDs
	var chatroomIDs []string
	cursor, err := r.chatroomCollection.Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"_id": 1}))
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var chatroom struct {
			ID string `bson:"_id"`
		}
		if err := cursor.Decode(&chatroom); err != nil {
			return err
		}
		chatroomIDs = append(chatroomIDs, chatroom.ID)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	// Delete messages where chatroom_id is not in chatroomIDs
	_, err = r.messageCollection.DeleteMany(context.TODO(), bson.M{
		"chatroom_id": bson.M{"$nin": chatroomIDs},
	})
	return err
}

func GetUserChatroomPersonalConfig(r *MongoChatroomRepository, userID string, page, pageSize int) ([]*domain.ChatroomPersonalConfig, error) {
	// Validate page and pageSize parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // Default pageSize
	}

	userFilter := bson.M{"user_id": userID}

	// Calculate skip and limit
	skip := (page - 1) * pageSize
	limit := pageSize

	// Create aggregation pipeline
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: userFilter}},
		bson.D{{Key: "$unwind", Value: "$chatroom_list"}},
		bson.D{{Key: "$skip", Value: skip}},
		bson.D{{Key: "$limit", Value: limit}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "chatroom_list", Value: 1}}}},
	}

	// Execute the aggregation
	cursor, err := r.userChatroomCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(context.TODO())

	var personalConfigs []*domain.ChatroomPersonalConfig
	for cursor.Next(context.TODO()) {
		var result struct {
			ChatroomList domain.ChatroomPersonalConfig `bson:"chatroom_list"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode chatroom config: %w", err)
		}
		personalConfigs = append(personalConfigs, &result.ChatroomList)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return personalConfigs, nil
}

func (r *MongoChatroomRepository) RemoveUserFromChatroom(obj domain.SpecifyChatRoomUser) error {
	// Step 1: Remove user from chatroomCollection
	chatroomFilter := bson.M{
		"_id": obj.ChatroomID,
	}

	chatroomUpdate := bson.M{
		"$pull": bson.M{
			"user_list": obj.UserID,
		},
	}

	result, err := r.chatroomCollection.UpdateOne(context.TODO(), chatroomFilter, chatroomUpdate)
	if err != nil {
		return fmt.Errorf("failed to remove user from chatroom: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("chatroom does not exist or user is not in chatroom")
	}

	// Step 2: Remove chatroom from userChatroomCollection
	userChatroomFilter := bson.M{
		"user_id":                   obj.UserID,
		"chatroom_list.chatroom_id": obj.ChatroomID,
	}

	userChatroomUpdate := bson.M{
		"$pull": bson.M{
			"chatroom_list": bson.M{
				"chatroom_id": obj.ChatroomID,
			},
		},
	}

	result, err = r.userChatroomCollection.UpdateOne(context.TODO(), userChatroomFilter, userChatroomUpdate)
	if err != nil {
		return fmt.Errorf("failed to update user chatroom list: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("user chatroom record does not exist or chatroom is not in user's list")
	}

	// Step 3: Optionally remove empty chatrooms
	err = r.RemoveEmptyChatrooms()
	if err != nil {
		return fmt.Errorf("error removing empty chatrooms: %w", err)
	}

	return nil
}

func (r *MongoChatroomRepository) RemoveEmptyChatrooms() error {
	// Define a filter to find chatrooms where user_list array is empty
	filter := bson.M{"user_list": bson.M{"$size": 0}}

	// Perform the delete operation
	_, err := r.chatroomCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to remove empty chatrooms: %w", err)
	}

	return nil
}

// UpdateUserChatroom updates the chatroom list of a user in userChatroomCollection
func (r *MongoChatroomRepository) UpdateUserChatroomList(userID string, chatroomList []domain.ChatroomPersonalConfig) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"chatroom_list": chatroomList}}

	_, err := r.userChatroomCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r *MongoChatroomRepository) UpdateChatroomUserList(obj domain.SpecifyChatRoomUser) error {
	// Filter to find the chatroom by ID
	filter := bson.M{"_id": obj.ChatroomID}

	// Update to add the userID to the user_list array if it doesn't already exist
	update := bson.M{
		"$addToSet": bson.M{
			"user_list": obj.UserID,
		},
	}

	// Perform the update operation
	_, err := r.chatroomCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update chatroom user list: %w", err)
	}

	return nil
}

func (r *MongoChatroomRepository) GetUserChatroomByUserID(userID string, chatroomID string) (*domain.ChatroomPersonalConfig, error) {
	filter := bson.M{"user_id": userID}

	var userChatroom domain.Entity_UserChatroom
	err := r.userChatroomCollection.FindOne(context.TODO(), filter).Decode(&userChatroom)
	if err == mongo.ErrNoDocuments {
		return nil, nil // Return nil if the user does not have any chatroom records
	}
	if err != nil {
		return nil, err
	}

	// Iterate through the ChatroomList to find the chatroom with the given ChatroomID
	for _, chatroom := range userChatroom.ChatroomList {
		if chatroom.ChatroomID == chatroomID {
			// Return the specific chatroom if found
			return &chatroom, nil
		}
	}

	// Return nil if the specific chatroomID is not found in the list
	return nil, nil
}

func (r *MongoChatroomRepository) GetChatroomUserList(chatroomID string) ([]string, error) {
	filter := bson.M{"_id": chatroomID}

	var result struct {
		UserList []string `bson:"user_list"`
	}

	err := r.chatroomCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("chatroom not found")
		}
		return nil, fmt.Errorf("failed to retrieve chatroom user list: %w", err)
	}

	return result.UserList, nil
}

func (r *MongoChatroomRepository) SetUserChatRoomConfig(dto domain.SetUserChatRoomConfig) error {
	// Step 1: Find the user chatroom record
	filter := bson.M{"user_id": dto.UserID}
	update := bson.M{
		"$set": bson.M{
			"chatroom_list.$[elem]." + dto.ParamName: dto.ParamValue,
		},
	}

	// Update options to specify array filter
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.chatroom_id": dto.ChatroomID},
		},
	})

	// Perform the update operation
	_, err := r.userChatroomCollection.UpdateOne(context.TODO(), filter, update, arrayFilters)
	if err != nil {
		return fmt.Errorf("failed to update user chatroom config: %w", err)
	}

	return nil
}

func (r *MongoChatroomRepository) AddUserToChatroom(obj domain.SpecifyChatRoomUser, personalConfig *domain.ChatroomPersonalConfig) error {
	// Start a session
	session, err := r.client.StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(context.Background())

	// Start a transaction
	_, err = session.WithTransaction(context.Background(), func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Step 1: Add userID to chatroomCollection user_list
		chatroomFilter := bson.M{"_id": obj.ChatroomID}
		chatroomUpdate := bson.M{
			"$addToSet": bson.M{
				"user_list": obj.UserID,
			},
		}
		_, err := r.chatroomCollection.UpdateOne(sessCtx, chatroomFilter, chatroomUpdate)
		if err != nil {
			return nil, fmt.Errorf("failed to update chatroom user list: %w", err)
		}

		// Step 2: Update user's chatroom list
		userChatroomFilter := bson.M{"user_id": obj.UserID}
		userChatroomUpdate := bson.M{
			"$addToSet": bson.M{
				"chatroom_list": personalConfig,
			},
		}

		result, err := r.userChatroomCollection.UpdateOne(
			sessCtx,
			userChatroomFilter,
			userChatroomUpdate,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update user chatroom list: %w", err)
		}

		// If no document was modified and no new document was inserted, it means the chatroom already exists
		if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
			return nil, fmt.Errorf("user is already in the chatroom")
		}

		return nil, nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}

func (r *MongoChatroomRepository) GetUserChatrooms(userID string, page, pageSize int) ([]domain.Entity_ChatroomsList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1000 {
		pageSize = 1000 // Default pageSize
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{
		// Stage 1: Match documents
		bson.D{{
			Key: "$match",
			Value: bson.D{{
				Key:   "user_id",
				Value: userID,
			}},
		}},
		// Stage 2: Unwind the chatroom_list array
		bson.D{{
			Key:   "$unwind",
			Value: "$chatroom_list",
		}},
		// Stage 3: Project chatroom_id
		bson.D{{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "chatroom_id", Value: "$chatroom_list.chatroom_id"},
			},
		}},
		// Stage 4: Lookup chatrooms data
		bson.D{{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "chatrooms"},         // Target collection name
				{Key: "localField", Value: "chatroom_id"}, // Field in the current collection
				{Key: "foreignField", Value: "_id"},       // Field in the target collection
				{Key: "as", Value: "chatroom_data"},       // Output field name
			},
		}},
		// Stage 5: Unwind the chatroom_data array
		bson.D{{
			Key:   "$unwind",
			Value: "$chatroom_data",
		}},
		// Stage 6: Project chatroom_data
		bson.D{{
			Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "chatroom_data", Value: 1},
			},
		}},
	}

	// Execute the aggregation
	cursor, err := r.userChatroomCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(context.Background())

	// Check cursor error
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	var results []domain.Entity_ChatroomsList

	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("failed to decode cursor document: %w", err)
		}

		// Safely assert and extract fields
		chatroomData, ok := document["chatroom_data"].(bson.M)
		if !ok {
			return nil, fmt.Errorf("invalid chatroom_data type")
		}

		// Initialize the domain struct with data
		entityChatroom := domain.Entity_Chatroom{
			ID:            getField(chatroomData, "_id", ""),
			BaseID:        getField(chatroomData, "base_id", ""),
			BoardID:       getField(chatroomData, "board_id", ""),
			Title:         getField(chatroomData, "title", ""),
			LastMessageID: getField(chatroomData, "last_message_id", ""),
			Avatar:        getField(chatroomData, "avatar", ""),
			CreatorID:     getField(chatroomData, "creator_id", ""),
			CreatedAt:     getDateTimeField(chatroomData, "created_at"),
			UpdatedAt:     getDateTimeField(chatroomData, "updated_at"),
			Type:          getField(chatroomData, "type", ""),
		}
		lastMessageID := getFieldValue(chatroomData, "last_message_id", "")
		lastMessage, err := r.findLastMessage(context.Background(), lastMessageID.(string))
		if err != nil {
			return nil, fmt.Errorf("failed to find last message: %v", err)
		}
		// Convert to your desired interface{}
		chatroomsLists := &domain.Entity_ChatroomsList{
			Chatroom: entityChatroom,
			PersonalConfig: domain.ChatroomPersonalConfig{
				ChatroomID:   getField(chatroomData, "chatroom_id", ""),
				LastReadTime: getDateTimeField(chatroomData, "last_read_time"),
				UnreadCount:  getField(chatroomData, "unread_count", 0),
				IsMuted:      getField(chatroomData, "is_muted", false),
				IsHidden:     getField(chatroomData, "is_hidden", false),
				IsPinned:     getField(chatroomData, "is_pinned", false),
			},
			LastMessage: domain.LastMessage{
				Content:      lastMessage.Content,
				Timestamp:    lastMessage.Timestamp,
				SenderID:     lastMessage.SenderID,
				SenderName:   lastMessage.SenderName,
				SenderAvatar: lastMessage.SenderAvatar,
			},
		}

		results = append(results, *chatroomsLists)
	}

	return results, nil
}

func (r *MongoChatroomRepository) findLastMessage(ctx context.Context, lastMessageID string) (*domain.LastMessage, error) {
	// objectID, err := primitive.ObjectIDFromHex(lastMessageID)
	// if err != nil {
	// 	return &domain.LastMessage{
	// 		Content:      "",
	// 		Timestamp:    time.Time{},
	// 		SenderID:     "",
	// 		SenderName:   "",
	// 		SenderAvatar: "",
	// 	}, nil
	// }
	var entityMessage domain.Entity_Message

	filter := bson.M{"_id": lastMessageID}

	err := r.messageCollection.FindOne(ctx, filter).Decode(entityMessage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &domain.LastMessage{
				Content:      "",
				Timestamp:    time.Time{},
				SenderID:     "",
				SenderName:   "",
				SenderAvatar: "",
			}, nil
		}
		return nil, fmt.Errorf("failed to find last message: %v", err)
	}
	lastMessage := &domain.LastMessage{
		Content:      entityMessage.Data,
		Timestamp:    entityMessage.UpdatedAt,
		SenderID:     entityMessage.SenderID,
		SenderName:   entityMessage.Nickname,
		SenderAvatar: entityMessage.Avatar,
	}
	return lastMessage, nil
}

func getDateTimeField(data bson.M, key string) time.Time {
	if value, ok := data[key].(primitive.DateTime); ok {
		return value.Time()
	}
	return time.Time{}
}
func getFieldValue(doc bson.M, key string, defaultValue interface{}) interface{} {
	if value, exists := doc[key]; exists {
		if reflect.TypeOf(value) == reflect.TypeOf(defaultValue) {
			return value
		}
	}
	return defaultValue
}

func getField[T any](doc bson.M, key string, defaultValue T) T {
	if value, ok := doc[key].(T); ok {
		return value
	}
	return defaultValue
}
