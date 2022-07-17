package users

import (
	"fmt"
	"sync"
	"testing"

	storagedb "posterr/src/storage/db"
	"posterr/src/storage/posterr"
	testdb "posterr/src/test/db"
	testrand "posterr/src/test/rand"

	assertions "github.com/stretchr/testify/assert"
)

const (
	maxContentSize    = 777
	maxUsernameLength = 14
)

func TestUserCreation(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	t.Run("Many random names", func(t *testing.T) {
		for count := 0; count < 100; count++ {
			username := rs.GenerateUnique(maxUsernameLength)
			err = users.CreateUser(username)
			assert.NoError(err)
		}
	})

	t.Run("Username too big", func(t *testing.T) {
		username := rs.GenerateUnique(maxUsernameLength + 1)
		err = users.CreateUser(username)
		assert.Equal(UsernameExceededMaximumCharsError{username}, err)
	})

	t.Run("Username with invalid characters", func(t *testing.T) {
		username := fmt.Sprintf("%s@", rs.GenerateUnique(maxUsernameLength-1))
		err = users.CreateUser(username)
		assert.Equal(InvalidUsernameError{username}, err)
	})
}

func TestFollowUser(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	t.Run("Parallel follows", func(t *testing.T) {
		// TODO: test fails for too many connections
		noUsers := 10
		usernames := make([]string, noUsers, noUsers)
		for i := 0; i < noUsers; i++ {
			usernames[i] = rs.GenerateUnique(maxUsernameLength)
			err = users.CreateUser(usernames[i])
			assert.NoError(err)
		}

		var wg sync.WaitGroup
		follow := func(userA, userB string, wg *sync.WaitGroup) {
			err := users.FollowUser(userA, userB)
			assert.NoError(err)
			wg.Done()
		}

		for a := 0; a < noUsers; a++ {
			for b := a + 1; b < noUsers; b++ {
				wg.Add(1)
				go follow(usernames[a], usernames[b], &wg)
			}
		}
		wg.Wait()
	})

	t.Run("Should not follow itself", func(t *testing.T) {
		username := rs.GenerateUnique(maxUsernameLength)
		err = users.FollowUser(username, username)
		assert.Error(err)
	})

	t.Run("Should fail if user already follows", func(t *testing.T) {
		userA := rs.GenerateUnique(maxUsernameLength)
		userB := rs.GenerateUnique(maxUsernameLength)
		err = users.CreateUser(userA)
		assert.NoError(err)
		err = users.CreateUser(userB)
		assert.NoError(err)

		err = users.FollowUser(userA, userB)
		assert.NoError(err)

		err = users.FollowUser(userA, userB)
		assert.Error(err)
	})
}

func TestUnfollowUser(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	t.Run("Parallel unfollows", func(t *testing.T) {
		// TODO: test fails for too many connections
		noUsers := 10
		usernames := make([]string, noUsers, noUsers)
		for i := 0; i < noUsers; i++ {
			usernames[i] = rs.GenerateUnique(maxUsernameLength)
			err = users.CreateUser(usernames[i])
			assert.NoError(err)
		}

		var wg sync.WaitGroup
		follow := func(userA, userB string, wg *sync.WaitGroup) {
			err := users.FollowUser(userA, userB)
			assert.NoError(err)
			wg.Done()
		}

		for a := 0; a < noUsers; a++ {
			for b := a + 1; b < noUsers; b++ {
				wg.Add(1)
				go follow(usernames[a], usernames[b], &wg)
			}
		}
		wg.Wait()

		// after creating the list of followers, try to unfollow each pair
		unfollow := func(userA, userB string, wg *sync.WaitGroup) {
			err := users.UnfollowUser(userA, userB)
			assert.NoError(err)
			wg.Done()
		}

		for a := 0; a < noUsers; a++ {
			for b := a + 1; b < noUsers; b++ {
				wg.Add(1)
				go unfollow(usernames[a], usernames[b], &wg)
			}
		}
		wg.Wait()
	})

	t.Run("Should not unfollow itself", func(t *testing.T) {
		username := rs.GenerateUnique(maxUsernameLength)
		err = users.UnfollowUser(username, username)
		assert.Error(err)
	})

	t.Run("Should fail if user does not follow", func(t *testing.T) {
		userA := rs.GenerateUnique(maxUsernameLength)
		userB := rs.GenerateUnique(maxUsernameLength)
		// if users do not exist, it will fail as well
		err = users.UnfollowUser(userA, userB)
		assert.Error(err)
	})
}

func TestIsFollowingUser(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	userA := rs.GenerateUnique(maxUsernameLength)
	userB := rs.GenerateUnique(maxUsernameLength)
	err = users.CreateUser(userA)
	assert.NoError(err)
	err = users.CreateUser(userB)
	assert.NoError(err)

	err = users.FollowUser(userA, userB)
	assert.NoError(err)

	t.Run("Should be true if userA follows userB", func(t *testing.T) {
		value, err := users.IsFollowingUser(userA, userB)
		assert.NoError(err)
		assert.True(value)
	})

	t.Run("Should be false if userB does not follow userA", func(t *testing.T) {
		value, err := users.IsFollowingUser(userB, userA)
		assert.NoError(err)
		assert.False(value)
	})
}

func TestCountUserPosts(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	posts := posterr.NewPosterrBacked(db)
	users := NewUserBacked(db)

	username := rs.GenerateUnique(14)
	err = users.CreateUser(username)
	assert.NoError(err)

	noPosts := 4
	for i := 0; i < noPosts; i++ {
		content := rs.GenerateAny(maxContentSize)
		_, err = posts.WriteContent(username, content, "")
		assert.NoError(err)
	}

	t.Run("Should match no. of posts for a username", func(t *testing.T) {
		count, err := users.CountUserPosts(username)
		assert.NoError(err)
		assert.Equal(noPosts, count)
	})

	t.Run("Should be empty if user does not have posts", func(t *testing.T) {
		count, err := users.CountUserPosts("noSuchUser")
		assert.NoError(err)
		assert.Empty(count)
	})
}

func TestUserFollowers(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	noUsers := 10
	usernames := make([]string, noUsers, noUsers)
	for i := 0; i < noUsers; i++ {
		usernames[i] = rs.GenerateUnique(maxUsernameLength)
		err = users.CreateUser(usernames[i])
		assert.NoError(err)
	}

	for a := 0; a < noUsers-1; a++ {
		err := users.FollowUser(usernames[0], usernames[a+1])
		assert.NoError(err)
	}

	count, err := users.CountUserFollowers(usernames[0])
	assert.NoError(err)
	assert.Equal(noUsers-1, count)
}

func TestUserFollowing(t *testing.T) {
	assert := assertions.New(t)
	rs := testrand.NewPseudoRandomString()
	dbName := testdb.GenerateDBName()

	db := storagedb.NewDatabase(dbName)
	err := db.InitializeDB()
	defer testdb.DropDatabase(dbName)
	assert.NoError(err)

	users := NewUserBacked(db)

	noUsers := 10
	usernames := make([]string, noUsers, noUsers)
	for i := 0; i < noUsers; i++ {
		usernames[i] = rs.GenerateUnique(maxUsernameLength)
		err = users.CreateUser(usernames[i])
		assert.NoError(err)
	}

	for a := 0; a < noUsers-1; a++ {
		err := users.FollowUser(usernames[a+1], usernames[0])
		assert.NoError(err)
	}

	count, err := users.CountUserFollowing(usernames[0])
	assert.NoError(err)
	assert.Equal(noUsers-1, count)
}
