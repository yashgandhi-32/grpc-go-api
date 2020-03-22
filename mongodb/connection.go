package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/yashgandhi-32/GRPC-API-CRUD/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongotest = "mongodb://localhost:27017"
)

type ConnectionManager struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func ConnectDB() (*errors.Message, *ConnectionManager) {
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongotest))
	if err != nil {
		log.Fatalf("Failed to connect mongodb %v", err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		return errors.Wrap(err, "An error connection DB"), nil
	}
	db := client.Database("blog")
	return nil, &ConnectionManager{Client: client, Db: db}
}
