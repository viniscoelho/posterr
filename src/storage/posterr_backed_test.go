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
		err = posts.WritePost(username, "", 0)
		assert.Error(err)
	})

	t.Run("Content just right", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		err = posts.WritePost(username, content, 0)
		assert.NoError(err)
	})

	t.Run("Content too long", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize + 1)
		err = posts.WritePost(username, content, 0)
		assert.Error(err)
	})
}

func TestTooManyPosts(t *testing.T) {
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
		err = posts.WritePost(username, content, 0)
		assert.NoError(err)
	}

	content := rs.GenerateAny(maxContentSize)
	err = posts.WritePost(username, content, 0)
	assert.Error(err)
}

// TODO: might be necessary to update write post to return the id of the post
