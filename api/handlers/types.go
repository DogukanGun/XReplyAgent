package handlers

import (
	"cg-mentions-bot/internal/utils/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ContextKey string

const UidKey ContextKey = "uid"

// Database connection
var DBClient *mongo.Client
var UserCollection = db.MongoDB{
	Database:   "xreplyagent",
	Collection: "users",
}

// User model for database
type User struct {
	FirebaseID string `bson:"firebase_id" json:"firebase_id"`
	TwitterID  string `bson:"twitter_id" json:"twitter_id"`
	Username   string `bson:"username" json:"username"`
}

type CheckUserResponse struct {
	Register bool `json:"register"`
}

type RegisterUserRequest struct {
	TwitterID string `json:"twitter_id"`
	Username  string `json:"username"`
}

type RegisterUserResponse struct {
	UID       string `json:"uid"`
	TwitterID string `json:"twitter_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
}

type ExecuteAppRequest struct {
	AppName string                 `json:"app_name"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type ExecuteAppResponse struct {
	UID     string `json:"uid"`
	AppName string `json:"app_name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ProfileResponse struct {
	UID       string `json:"uid"`
	Username  string `json:"username,omitempty"`
	TwitterID string `json:"twitter_id,omitempty"`
	Message   string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
