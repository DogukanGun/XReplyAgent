package db

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

func (mg MongoDB) Update(client *mongo.Client, filter interface{}, update interface{}) bool {
	if client == nil {
		log.Printf("Error in update function: mongo client is nil")
		return false
	}

	collection := client.Database(mg.Database).Collection(mg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error in update function: %s", err)
		return false
	}
	if res == nil {
		log.Printf("Error in update function: update result is nil")
		return false
	}

	return res.ModifiedCount > 0
}
