package gmo

import (
	"context"
	"fmt"
	"testing"

	"github.com/npalladium/gmo/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDB(t *testing.T) {
	ctx := context.Background()

	container, err := internal.StartContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	// check that mongodb is working
	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Error(fmt.Errorf("failed to get endpoint: %w", err))
	}
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatal(fmt.Errorf("error creating mongo client: %w", err))
	}
	err = mongoClient.Connect(ctx)
	if err != nil {
		t.Fatal(fmt.Errorf("error connecting to mongo: %w", err))
	}
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error pinging mongo: %w", err))
	}

	connString := endpoint + "/test"
	db, err := New(connString)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		t.Fatal(err)
	}

}
