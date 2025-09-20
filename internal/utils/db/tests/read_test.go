package tests

import (
	"cg-mentions-bot/internal/utils/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestRead(t *testing.T) {
	t.Run("Any object in database can be read", func(t *testing.T) {
		client, err := db.ConnectToDB("mongodb://localhost:27017/")
		if err != nil {
			t.Failed()
		}
		mg := db.MongoDB{
			Database:   "xreply_test",
			Collection: "db_test",
		}
		var dbItems []interface{}
		mg.Read(client, bson.D{}, &dbItems)
		assert.Less(t, 0, len(dbItems))
	})
}
