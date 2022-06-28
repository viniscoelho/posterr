package storage

import (
	"context"
	"fmt"
	"log"
	typesrand "posterr/src/types/rand"
	"sync"
	"testing"

	"posterr/src/storage/postgres"

	assertions "github.com/stretchr/testify/assert"
)

const maxUsernameLength = 14

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
		for count := 0; count < 100; count++ {
			username := rs.GenerateUnique(maxUsernameLength)
			err = users.CreateUser(username)
			assert.NoError(err)
		}
	})

	t.Run("Username too big", func(t *testing.T) {
		username := rs.GenerateUnique(maxUsernameLength + 1)
		err = users.CreateUser(username)
		assert.Error(err)
	})

	t.Run("Username with invalid characters", func(t *testing.T) {
		username := fmt.Sprintf("%s@", rs.GenerateUnique(maxUsernameLength-1))
		err = users.CreateUser(username)
		assert.Error(err)
	})
}

func TestFollowUser(t *testing.T) {
	assert := assertions.New(t)
	dbName := "dummy"

	db := postgres.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)
	rs := typesrand.NewPseudoRandomString()

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
	dbName := "dummy"

	db := postgres.NewDatabase(dbName)
	err := db.InitializeDB()
	defer dropDatabase(dbName)
	assert.NoError(err)

	posts := NewPosterrBacked(db)
	users := NewUserBacked(db, posts)
	rs := typesrand.NewPseudoRandomString()

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

func dropDatabase(dbName string) {
	db := postgres.NewDatabase("postgres")

	conn, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, dbName))
	if err != nil {
		log.Printf("Failed to drop database %s: %s", dbName, err)
		return
	}

	log.Printf("Database %s dropped!", dbName)
}
