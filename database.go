package gmo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Database struct {
	db     *mongo.Database
	client *mongo.Client
	DBName string
}

// New creates an instance of gmo.Database.
func New(connString string) (database Database, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	parsedConnString, err := connstring.ParseAndValidate(connString)
	if err != nil {
		return database, err
	}

	clientOptions := options.Client().ApplyURI(connString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return database, err
	}
	database.client = client

	err = database.client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Panic("Could not connect to MongoDB! Please check if mongo is running.", err)
		return database, err
	}
	log.Print("connected to MongoDB")

	database.DBName = parsedConnString.Database
	database.db = database.client.Database(database.DBName)
	return database, err
}

// ListCollectionNames lists all the collections in the DB.
func (d *Database) ListCollectionNames(ctx context.Context, filter any, opts ...*options.ListCollectionsOptions) ([]string, error) {
	return d.db.ListCollectionNames(ctx, filter, opts...)
}

// DefaultContext creates a context for convenience.
func DefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}
