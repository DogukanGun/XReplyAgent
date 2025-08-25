package db

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

func (mg MongoDB) Insert(client *mongo.Client, obj interface{}) bool {
	collection := client.Database(mg.Database).Collection(mg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, obj)
	if err != nil {
		log.Printf("Error in insert function: %s", err)
	}
	return res.Acknowledged
}
