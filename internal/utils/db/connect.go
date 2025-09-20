package db

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectToDB(connectionUri string) (*mongo.Client, error) {
	client, _ := mongo.Connect(options.Client().ApplyURI(connectionUri))
	return client, nil
}
