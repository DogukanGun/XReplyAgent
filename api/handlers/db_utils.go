package handlers

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SetDBClient sets the database client for handlers
func SetDBClient(client *mongo.Client) {
	DBClient = client
}

// CheckUserExists checks if user exists in database by Firebase ID
func CheckUserExists(firebaseID string) bool {
	if DBClient == nil {
		return false
	}

	filter := bson.D{{"firebase_id", firebaseID}}
	var users []User

	exists := UserCollection.Read(DBClient, filter, &users)
	return exists && len(users) > 0
}

// CreateUser creates a new user in database
func CreateUser(firebaseID, twitterID, username string) bool {
	if DBClient == nil {
		return false
	}

	user := User{
		FirebaseID: firebaseID,
		TwitterID:  twitterID,
		Username:   username,
	}

	return UserCollection.Insert(DBClient, user)
}

// GetUserByFirebaseID gets user from database by Firebase ID
func GetUserByFirebaseID(firebaseID string) (*User, bool) {
	if DBClient == nil {
		return nil, false
	}

	filter := bson.D{{"firebase_id", firebaseID}}
	var users []User

	if UserCollection.Read(DBClient, filter, &users) && len(users) > 0 {
		return &users[0], true
	}

	return nil, false
}
