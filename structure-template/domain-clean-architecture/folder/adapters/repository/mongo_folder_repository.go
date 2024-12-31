package repository

import (
	"context"
	"errors"
	"fmt"
	"folder_mod/core/domain/folder"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFolderRepository struct {
	collection *mongo.Collection
}

func NewMongoFolderRepository(client *mongo.Client) folder.FolderRepository {
	collection := client.Database("mai_dev").Collection("folders")
	return &MongoFolderRepository{collection: collection}
}

func (r *MongoFolderRepository) CreateFolder(folder *folder.Folder) (*folder.Folder, error) {
	newID := primitive.NewObjectID()
	folder.ID = newID.Hex()
	result, err := r.collection.InsertOne(context.TODO(), folder)

	if err != nil {
		fmt.Printf("Error inserting document: %v\n", err)
		return nil, err
	}
	fmt.Printf("result: %v\n", result)
	return folder, nil
}

func (r *MongoFolderRepository) GetFolders() ([]*folder.Folder, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var folders []*folder.Folder
	if err := cursor.All(context.TODO(), &folders); err != nil {
		return nil, err
	}
	return folders, nil
}

// func (r *MongoFolderRepository) FindByID(id string) (*folder.Folder, error) {
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var folder folder.Folder
// 	if err := r.collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&folder); err != nil {
// 		return nil, err
// 	}
// 	return &folder, nil
// }

func (r *MongoFolderRepository) UpdateFolderData(_folder *folder.Folder) (*folder.Folder, error) {
	// Convert folder ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(_folder.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid folder ID: %w", err)
	}

	// Define filter and update operations
	filter := bson.M{"_id": objectID, "base_id": _folder.BaseID}
	update := bson.M{
		"$set": bson.M{
			"parent_id":  _folder.ParentID,
			"position":   _folder.Position,
			"data":       _folder.Data,
			"updated_at": _folder.UpdatedAt,
		},
	}

	// Perform the update operation
	result, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update folder: %w", err)
	}

	// Check if a document was matched and updated
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no matching document found for ID: %s and BaseID: %s", _folder.ID, _folder.BaseID)
	}

	// Retrieve and return the updated folder
	var updatedFolder folder.Folder
	err = r.collection.FindOne(context.TODO(), filter).Decode(&updatedFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated folder: %w", err)
	}

	return &updatedFolder, nil
}

func (r *MongoFolderRepository) DeleteFolderByID(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %w", err)
	}

	// Find the folder to be deleted based on the _id
	var folderToDelete folder.Folder
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&folderToDelete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("folder not found")
		}
		return fmt.Errorf("failed to find folder: %w", err)
	}

	// Check if the folder type is "folder"
	if folderToDelete.Type != "folder" {
		return fmt.Errorf("item to delete is not a folder")
	}

	// Recursively delete folder and its children
	return r.deleteFolderAndChildren(objectID)
}

func (r *MongoFolderRepository) DeleteFolder(_folder *folder.Folder) error {
	objectID, err := primitive.ObjectIDFromHex(_folder.ID)
	if err != nil {
		return fmt.Errorf("invalid parent_id: %w", err)
	}
	// Find the folder to be deleted based on the _id
	var folderToDelete folder.Folder
	err = r.collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&folderToDelete)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("folder not found")
		}
		return fmt.Errorf("failed to find folder: %w", err)
	}

	// Recursively delete folder and its children
	return r.deleteFolderAndChildren(objectID)
}

func (r *MongoFolderRepository) deleteFolderAndChildren(folderID primitive.ObjectID) error {
	// Find all children folders
	children, err := r.FindFoldersByParentID(folderID.Hex())
	if err != nil {
		return err
	}

	// Recursively delete all children folders
	for _, child := range children {
		childID, err := primitive.ObjectIDFromHex(child.ID)
		if err != nil {
			return fmt.Errorf("invalid child ID: %w", err)
		}
		if err := r.deleteFolderAndChildren(childID); err != nil {
			return err
		}
	}

	// Delete the folder itself
	_, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": folderID})
	return err
}

func (r *MongoFolderRepository) FindFolderByObjectID(objectID string) (*folder.Folder, error) {
	// Convert the string ID to ObjectID
	id, err := primitive.ObjectIDFromHex(objectID)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID: %w", err)
	}

	// Define filter for querying
	filter := bson.M{"_id": id}

	// Find the folder by ObjectID
	var folder folder.Folder
	err = r.collection.FindOne(context.TODO(), filter).Decode(&folder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no folder found with ID: %s", objectID)
		}
		return nil, fmt.Errorf("failed to find folder: %w", err)
	}

	return &folder, nil
}

func (r *MongoFolderRepository) FindFoldersByParentID(parentID string) ([]folder.Folder, error) {
	filter := bson.M{"parent_id": parentID}
	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var folders []folder.Folder
	for cursor.Next(context.TODO()) {
		var folder folder.Folder
		if err := cursor.Decode(&folder); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return folders, nil
}

func (r *MongoFolderRepository) UpdateFolderParentID(objectID string, parentID string) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"base_id": objectID},
		bson.M{"$set": bson.M{"parent_id": parentID}},
	)
	return err
}

func (r *MongoFolderRepository) AddChildIDToParent(parentID string, childID string) error {
	_, err := r.collection.UpdateOne(
		context.TODO(), // Using a default context
		bson.M{"base_id": parentID},
		bson.M{"$addToSet": bson.M{"child_ids": childID}},
	)
	return err
}

func (r *MongoFolderRepository) PositionExists(baseID string, parentID string, position float64) error {
	filter := bson.M{
		"base_id":   baseID,
		"position":  position,
		"parent_id": parentID,
	}
	count, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("position already exists for the same baseID and parentID")
	}
	return nil
}
