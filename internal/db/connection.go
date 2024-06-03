package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "customer"
)

var DB *mongo.Database

type DataBase struct {
	CurrentConn *mongo.Client
}

func Connect() *mongo.Client {
	dbName, err := getConnection()
	if err != nil {
		return &mongo.Client{}
	}
	return dbName
}

func getConnection() (client *mongo.Client, ctx context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("authservice")
	fmt.Println("Connected to MongoDB!")

	return client, ctx
}
