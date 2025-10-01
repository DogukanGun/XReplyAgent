package db

import (
	"context"
	"crypto/tls"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

func ConnectToDB(connectionUri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Configure TLS for MongoDB Atlas
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	clientOptions := options.Client().
		ApplyURI(connectionUri).
		SetTLSConfig(tlsConfig).
		SetServerSelectionTimeout(10 * time.Second).
		SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
