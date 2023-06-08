package gmo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://github.com/mongodb/mongo-go-driver/blob/master/mongo/collection.go

//var ErrNoDocuments = mongo.ErrNoDocuments

// Collection is a generic wrapper around mongo.Collection.
type Collection[T any] struct {
	collection *mongo.Collection
}

// CountDocuments returns the number of documents in the collection.
func (c *Collection[T]) CountDocuments(filter any) (int64, error) {
	count, err := c.collection.CountDocuments(DefaultContext(), filter)
	return count, err
}

// FindOne executes a find command and a document in the collection.
func (c *Collection[T]) FindOne(ctx context.Context, filter any,
	opts ...*options.FindOneOptions) (result T, err error) {
	singleResult := c.collection.FindOne(ctx, filter)
	if singleResult.Err() != nil {
		return result, singleResult.Err()
	}
	err = c.collection.FindOne(ctx, filter).Decode(&result)
	return result, err
}

// FindByID executes a find command based on the ID.
func (c *Collection[T]) FindByID(ctx context.Context, id string) (result T, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}
	return c.FindOne(ctx, bson.M{"_id": objectID})
}

// Find executes a find command and returns a slice of the matching documents.
func (c *Collection[T]) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (result []T, err error) {
	cursor, err := c.collection.Find(ctx, filter, opts...)
	if err != nil {
		return result, err
	}

	err = cursor.All(ctx, &result)

	return result, nil
}

// Insert executes an insert command to insert a single document into the collection.
func (c *Collection[T]) Insert(document T) (T, error) {
	_, err := c.collection.InsertOne(DefaultContext(), document)
	return document, err
}

// UpdateOne executes an update command to update at most one document in the collection.
func (c *Collection[T]) UpdateOne(ctx context.Context, filter any, document T) error {
	_, err := c.collection.UpdateOne(ctx, filter, bson.M{"$set": document})
	return err
}

// UpdateByID executes an update command based on the ID.
func (c *Collection[T]) UpdateByID(ctx context.Context, id string, document T) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return c.UpdateOne(ctx, bson.M{"_id": objectID}, document)
}
