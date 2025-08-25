package db

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"time"
)

func (mg MongoDB) Read(client *mongo.Client, filter bson.D, result interface{}) bool {
	collection := client.Database(mg.Database).Collection(mg.Collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Print(err)
		return false
	}
	defer cur.Close(ctx)

	// Use cursor.All to fill the slice directly
	if err := cur.All(ctx, result); err != nil {
		log.Print(err)
		return false
	}

	return true
}
