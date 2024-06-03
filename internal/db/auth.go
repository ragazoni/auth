package db

import (
	"auth/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUserByEmail(collection string, email any) (*models.Customer, error) {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection("users")

	filter := bson.M{"email": email}

	var user models.Customer
	err := c.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	// Converte o ID para primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return nil, err
	}
	user.ID = id
	return &user, nil
}
