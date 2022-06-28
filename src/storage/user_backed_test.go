package storage

import (
	"context"
	"fmt"
	"log"
	typesrand "posterr/src/types/gen"
	"testing"

	"posterr/src/storage/postgres"

	assertions "github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	assert := assertions.New(t)
	dbName := "dummy"

	db := postgres.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)
	rs := typesrand.NewPseudoRandomString()

	t.Run("Many random names", func(t *testing.T) {
		for count := 0; count < 10; count++ {
			username := rs.Generate(14)
			err = users.CreateUser(username)
			assert.NoError(err)
		}
	})

	t.Run("Username too big", func(t *testing.T) {
		username := rs.Generate(15)
		err = users.CreateUser(username)
		assert.Error(err)
	})

	t.Run("Username with invalid characters", func(t *testing.T) {
		username := fmt.Sprintf("%s@", rs.Generate(13))
		err = users.CreateUser(username)
		assert.Error(err)
	})
}

func dropDatabase(dbName string) {
	db := postgres.NewDatabase("")

	conn, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, dbName))
	if err != nil {
		log.Printf("Failed to drop database %s: %s", dbName, err)
	}

	log.Printf("Database %s dropped!", dbName)
}
