package db

import (
	"auth/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Insert(collection string, data any) (primitive.ObjectID, error) {
	client, cxt := getConnection()
	defer mongo.Connect(cxt)

	c := client.Database(dbName).Collection(collection)

	resp, err := c.InsertOne(context.Background(), data)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return resp.InsertedID.(primitive.ObjectID), nil

}

func Find(collection string, document any) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection(collection)

	cursor, err := c.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	return cursor.All(context.Background(), document)
}

func FindById(collection, id string, document any) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection(collection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	return c.FindOne(context.Background(), filter).Decode(document)

}

func UpdateById(collection, id string, data, result any) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection(collection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID}
	opts := options.FindOneAndUpdate().SetUpsert(false)
	err = c.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": data}, opts).Err()
	if err != nil {
		return err
	}
	return c.FindOne(context.Background(), filter).Decode(result)

}

func DeletebyId(collection, id string) error {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection(collection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	result, err := c.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount != 1 {
		return fmt.Errorf("%d documents deleted", result.DeletedCount)
	}

	return nil
}

func FindByEmail(collection string, email any) (*models.Customer, error) {
	client, ctx := getConnection()
	defer client.Disconnect(ctx)

	c := client.Database(dbName).Collection(collection)

	filter := bson.M{"email": email}

	var result models.Customer
	err := c.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil

}
