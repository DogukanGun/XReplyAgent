package tests

import (
	"cg-mentions-bot/internal/utils/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestInsert(t *testing.T) {

	t.Run("Insert message can add a new valid object into db", func(t *testing.T) {
		client, err := db.ConnectToDB("mongodb://localhost:27017/")
		if err != nil {
			t.Fatal(err)
		}
		defer client.Disconnect(t.Context())

		mg := db.MongoDB{
			Database:   "xreply_test",
			Collection: "db_test",
		}

		if !mg.Insert(client, bson.D{{"test_id", "3435"}}) {
			t.Fatal("insert failed")
		}
	})

}
