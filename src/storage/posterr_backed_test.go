package storage

import (
	"testing"

	storagedb "posterr/src/storage/db"
	typesrand "posterr/src/types/rand"

	assertions "github.com/stretchr/testify/assert"
)

const maxContentSize = 777

func TestWritePost(t *testing.T) {
	assert := assertions.New(t)
	dbName := "dummy"

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)
	rs := typesrand.NewPseudoRandomString()

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Empty content", func(t *testing.T) {
		_, err = posts.WritePost(username, "", "")
		assert.Error(err)
	})

	t.Run("Content just right", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WritePost(username, content, "")
		assert.NoError(err)
	})

	t.Run("Content too long", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize + 1)
		_, err = posts.WritePost(username, content, "")
		assert.Error(err)
	})
}

func TestTooManyPostsInASingleDay(t *testing.T) {
	assert := assertions.New(t)
	dbName := "dummy"

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)
	rs := typesrand.NewPseudoRandomString()

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	for i := 0; i < 5; i++ {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WritePost(username, content, "")
		assert.NoError(err)
	}

	content := rs.GenerateAny(maxContentSize)
	_, err = posts.WritePost(username, content, "")
	assert.Error(err)
}

// TODO: might be necessary to update write post to return the id of the post
