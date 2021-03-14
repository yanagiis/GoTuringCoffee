package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BsonModel bson data which can be converted to lib object
type BsonModel interface {
	SetUpdatedTime()
}

// List list all items in the collection
func List(ctx context.Context, collection *mongo.Collection, result interface{}) (err error) {
	query := bson.M{}
	findOptions := options.Find()
	findOptions.Sort = bson.M{
		"created_at": -1,
	}

	rawResult, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return err
	}

	// allocate a value for the
	err = rawResult.All(ctx, result)
	if err != nil {
		return err
	}

	return nil
}

// GetByID query item by id
func GetByID(ctx context.Context, collection *mongo.Collection, id string, result interface{}) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}

// InsertOne insert a bson item to database
func InsertOne(ctx context.Context, collection *mongo.Collection, bsonItem interface{}) (oid primitive.ObjectID, err error) {
	newID, err := collection.InsertOne(ctx, bsonItem)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	oid = newID.InsertedID.(primitive.ObjectID)
	return
}

// ReplaceOne replace existing item
func ReplaceOne(ctx context.Context, collection *mongo.Collection, sid string, bsonItem interface{}) (err error) {
	objectID, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return err
	}

	bsonItem.(BsonModel).SetUpdatedTime()

	filter := bson.M{"_id": objectID}
	_, err = collection.ReplaceOne(ctx, filter, bsonItem)
	if err != nil {
		return err
	}

	return nil
}
