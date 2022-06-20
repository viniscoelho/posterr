package storage

import (
	"strings"
	"testing"

	"posterr/src/storage/postgres"

	assertions "github.com/stretchr/testify/assert"
)

const maxContentSize = 777

func TestWritePost(t *testing.T) {
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

	err = posts.WritePost("vinicius", "hello", 0)
	assert.NoError(err)

	err = posts.WritePost("vinicius", "", 0)
	assert.Error(err)

	justRight := strings.Repeat("1", maxContentSize)
	tooLong := strings.Repeat("1", maxContentSize+1)

	err = posts.WritePost("vinicius", justRight, 0)
	assert.NoError(err)

	err = posts.WritePost("vinicius", tooLong, 0)
	assert.Error(err)
}
