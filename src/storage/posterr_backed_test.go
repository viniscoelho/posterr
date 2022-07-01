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
	rs := typesrand.NewPseudoRandomString()
	dbName := generateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should fail for empty content", func(t *testing.T) {
		_, err = posts.WriteContent(username, "", "")
		assert.Error(err)
	})

	t.Run("Should post if content is just right", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WriteContent(username, content, "")
		assert.NoError(err)
	})

	t.Run("Should not post if content is too long", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize + 1)
		_, err = posts.WriteContent(username, content, "")
		assert.Error(err)
	})
}

func TestTooManyPostsInASingleDay(t *testing.T) {
	assert := assertions.New(t)
	rs := typesrand.NewPseudoRandomString()
	dbName := generateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	noPosts := 5
	for i := 0; i < noPosts; i++ {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WriteContent(username, content, "")
		assert.NoError(err)
	}

	content := rs.GenerateAny(maxContentSize)
	_, err = posts.WriteContent(username, content, "")
	assert.Error(err)
}

func TestRepost(t *testing.T) {
	assert := assertions.New(t)
	rs := typesrand.NewPseudoRandomString()
	dbName := generateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should repost an existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content, "")
		assert.NoError(err)

		_, err = posts.WriteContent(username, "", postId)
		assert.NoError(err)
	})

	t.Run("Should not repost an non existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err := posts.WriteContent(username, content, "")
		assert.NoError(err)

		_, err = posts.WriteContent(username, "", "somePostId")
		assert.Error(err)
	})
}

func TestQuotedRepost(t *testing.T) {
	assert := assertions.New(t)
	rs := typesrand.NewPseudoRandomString()
	dbName := generateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should quote repost an existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content, "")
		assert.NoError(err)

		_, err = posts.WriteContent(username, "check this out", postId)
		assert.NoError(err)
	})

	t.Run("Should not quote repost an non existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err := posts.WriteContent(username, content, "")
		assert.NoError(err)

		_, err = posts.WriteContent(username, "check this out", "somePostId")
		assert.Error(err)
	})
}
