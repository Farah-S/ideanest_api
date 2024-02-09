package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/golang-test/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	config, err := config.LoadDBConfig("./config/")
	
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx,options.Client().ApplyURI(config.DBUri))
	
	if err != nil {
		log.Panic(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()