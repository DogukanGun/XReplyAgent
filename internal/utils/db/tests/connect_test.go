package tests

import (
	"cg-mentions-bot/internal/utils/db"
	"testing"
)

func TestConnection(t *testing.T) {

	t.Run("Conncet function must have a connection with valid uri", func(t *testing.T) {
		client, err := db.ConnectToDB("mongodb://localhost:27017/")
		if err != nil {
			t.Failed()
			return
		}
		err = client.Ping(t.Context(), nil)
		if err != nil {
			t.Failed()
			return
		}
	})

	t.Run("Conncet function must fail with invalid uri", func(t *testing.T) {
		_, err := db.ConnectToDB("mongodb://localhost:2/")
		if err != nil {
			return
		}
		t.Failed()
	})
}
