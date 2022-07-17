package posterr

import (
	"testing"

	storagedb "posterr/src/storage/db"
	storageusers "posterr/src/storage/users"
	testdb "posterr/src/test/db"
	testrand "posterr/src/test/rand"

	assertions "github.com/stretchr/testify/assert"
)

const maxContentSize = 777

func TestWritePost(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := storageusers.NewUserBacked(db)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should post if content is just right", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WriteContent(username, content)
		assert.NoError(err)
	})

	t.Run("Should not post if content is too long", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize + 1)
		_, err = posts.WriteContent(username, content)
		assert.Equal(PostExceededMaximumCharsError{}, err)
	})
}

func TestTooManyPostsInASingleDay(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := storageusers.NewUserBacked(db)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	noPosts := 5
	for i := 0; i < noPosts; i++ {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WriteContent(username, content)
		assert.NoError(err)
	}

	content := rs.GenerateAny(maxContentSize)
	_, err = posts.WriteContent(username, content)
	assert.Equal(ExceededMaximumDailyPostsError{}, err)
}

func TestRepost(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := storageusers.NewUserBacked(db)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should repost an existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteRepostContent(username, postId)
		assert.NoError(err)
	})

	t.Run("Should not repost an non existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteRepostContent(username, "somePostId")
		assert.Equal(PostIdDoesNotExistError{"somePostId"}, err)
	})

	t.Run("Should not repost if username does not existing", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteRepostContent("notauser", postId)
		assert.Equal(UserDoesNotExistError{"notauser"}, err)
	})
}

func TestQuotedRepost(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := storageusers.NewUserBacked(db)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	t.Run("Should quote repost an existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteQuoteRepostContent(username, "check this out", postId)
		assert.NoError(err)
	})

	t.Run("Should not quote repost an non existing post", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		_, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteQuoteRepostContent(username, "check this out", "somePostId")
		assert.Equal(PostIdDoesNotExistError{"somePostId"}, err)
	})

	t.Run("Should not quote repost if username does not existing", func(t *testing.T) {
		content := rs.GenerateAny(maxContentSize)
		postId, err := posts.WriteContent(username, content)
		assert.NoError(err)

		_, err = posts.WriteQuoteRepostContent("notauser", "check this out", postId)
		assert.Equal(UserDoesNotExistError{"notauser"}, err)
	})
}
