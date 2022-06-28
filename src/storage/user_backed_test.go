package storage

import (
	"context"
	"fmt"
	"log"
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

	err = users.CreateUser("vinicius")
	assert.NoError(err)

	err = users.CreateUser("vinicius123456")
	assert.NoError(err)

	err = users.CreateUser("vinicius1234567")
	assert.Error(err)

	err = users.CreateUser("vinicius123456@")
	assert.Error(err)
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
