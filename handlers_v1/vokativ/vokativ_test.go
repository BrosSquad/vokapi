package vokativ_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/BrosSquad/vokapi/handlers_v1/vokativ"
	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func SetupDatabase(t testing.TB) *badger.DB {
	options := badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.ERROR)

	db, err := badger.Open(options)
	if err != nil {
		t.Errorf("Failed to create in memory database: %v", err)
		t.FailNow()
	}

	return db
}

var names = []string{
	"Dusan",
	"dusan",
	"Dušan",
	"dušan",
	"dUšAn",
	"DuŠAn",
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	app := fiber.New()
	db := SetupDatabase(t)

	app.Get("/:name", vokativ.Handler(db))

	for _, name := range names {
		t.Run(fmt.Sprintf("WithName_%s", name), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/Dusan", nil)

			res, err := app.Test(req)

			assert.NoError(err)
			assert.Equal(fiber.StatusNotFound, res.StatusCode)
		})
	}
}

func TestNameFound(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	app := fiber.New()
	db := SetupDatabase(t)

	tx := db.NewTransaction(true)

	_ = tx.Set([]byte("Dušan"), []byte("Dušane"))
	_ = tx.Set([]byte("Dusan"), []byte("Dušane"))
	_ = tx.Commit()


	app.Get("/:name", vokativ.Handler(db))

	for _, name := range names {
		t.Run(fmt.Sprintf("WithName_%s", name), func(t *testing.T) {
			req := httptest.NewRequest("GET", "/Dusan", nil)

			res, err := app.Test(req)

			assert.NoError(err)
			assert.Equal(fiber.StatusOK, res.StatusCode)
			m := map[string]string{}
			_ = json.NewDecoder(res.Body).Decode(&m)

			assert.Equal("Dušane", m["vokativ"])
		})
	}
}
