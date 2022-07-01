package storage

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	storagedb "posterr/src/storage/db"
	"posterr/src/types"
)

const (
	selectUser = `SELECT *
                 FROM users
                 WHERE username = $1`

	countFollowers = `SELECT COUNT(*) as followers
                 FROM followers
                 WHERE username = $1`

	countFollowing = `SELECT COUNT(*) as following
                 FROM followers
                 WHERE followed_by = $1`

	countUserPosts = `SELECT COUNT(*) as no_posts
                 FROM posts
                 WHERE username = $1`

	isFollowerOf = `SELECT COUNT(*) as is_follower
                 FROM followers
                 WHERE username = $1 AND followed_by = $2`
)

type FollowUser struct {
	Username string
	Follows  bool
}

type userBacked struct {
	sync.RWMutex
	// An accessor to the database
	db storagedb.ConnectDB
	// An accessor to the Posterr interface
	posts types.Posterr
	// A cache map to store how many followers a user has
	followersCount map[string]int
	// A cache map to store how many users a user is following
	followingCount map[string]int
	// A regex to validate usernames
	rgx *regexp.Regexp
}

func NewUserBacked(db storagedb.ConnectDB, posts types.Posterr) *userBacked {
	// TODO: should the connection pool be reset for every call?
	return &userBacked{
		db:             db,
		posts:          posts,
		followersCount: make(map[string]int),
		followingCount: make(map[string]int),
		rgx:            regexp.MustCompile(`^[a-zA-Z0-9]*$`),
	}
}

// CreateUser creates a user. This method is not exposed through an API
func (ub *userBacked) CreateUser(username string) error {
	if !ub.rgx.MatchString(username) {
		return fmt.Errorf("invalid username: a username must consist of alphanumeric charactes only")
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO users (username) VALUES ($1)", username)
	if err != nil {
		return fmt.Errorf("could not insert into users: %w", err)
	}

	return nil
}

// GetUserProfile returns a user detailed information
func (ub *userBacked) GetUserProfile(username string, offset int) (types.PosterrUser, error) {
	userProfile, err := ub.getUserDetails(username)
	if err != nil {
		return types.PosterrUser{}, err
	}

	userProfile.Followers, err = ub.CountUserFollowers(username)
	if err != nil {
		return types.PosterrUser{}, err
	}

	userProfile.Following, err = ub.CountUserFollowing(username)
	if err != nil {
		return types.PosterrUser{}, err
	}

	userProfile.PostsCount, err = ub.CountUserPosts(username)
	if err != nil {
		return types.PosterrUser{}, err
	}

	userProfile.Posts, err = ub.posts.ListProfileContent(username, offset)
	if err != nil {
		return types.PosterrUser{}, err
	}

	return userProfile, nil
}

// CountUserPosts returns how many posts a user has made
func (ub *userBacked) CountUserPosts(username string) (int, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countUserPosts, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform countUserPosts query: %w", err)
	}

	var dailyPosts int
	for rows.Next() {
		if err = rows.Scan(&dailyPosts); err != nil {
			return 0, fmt.Errorf("could not scan countUserPosts rows: %w", err)
		}
	}

	return dailyPosts, nil
}

// CountUserFollowing returns how many followers a user has
func (ub *userBacked) CountUserFollowers(username string) (int, error) {
	count, exists := func(username string) (int, bool) {
		ub.RLock()
		defer ub.RUnlock()
		count, exists := ub.followersCount[username]
		return count, exists
	}(username)

	if exists {
		return count, nil
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countFollowers, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform countFollowers query: %w", err)
	}

	var followers int
	for rows.Next() {
		if err = rows.Scan(&followers); err != nil {
			return 0, fmt.Errorf("could not scan countFollowers rows: %w", err)
		}
	}

	func(username string) {
		ub.Lock()
		defer ub.Unlock()
		ub.followersCount[username] = followers
	}(username)

	return followers, nil
}

// CountUserFollowing returns how many users a user is following
func (ub *userBacked) CountUserFollowing(username string) (int, error) {
	count, exists := func(username string) (int, bool) {
		ub.RLock()
		defer ub.RUnlock()
		count, exists := ub.followingCount[username]
		return count, exists
	}(username)

	if exists {
		return count, nil
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countFollowing, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform countFollowing query: %w", err)
	}

	var following int
	for rows.Next() {
		if err = rows.Scan(&following); err != nil {
			return 0, fmt.Errorf("could not scan countFollowing rows: %w", err)
		}
	}

	func(username string) {
		ub.Lock()
		defer ub.Unlock()
		ub.followingCount[username] = following
	}(username)

	return following, nil
}

// FollowUser ensures that username is followed by follower,
// i.e., follower follows username
func (ub *userBacked) FollowUser(username, follower string) error {
	if username == follower {
		return fmt.Errorf("%s cannot follow itself", username)
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	isFollowingUser, err := ub.IsFollowingUser(username, follower)
	if err != nil {
		return fmt.Errorf("could not check follower: %w", err)
	}

	if isFollowingUser {
		return UserAlreadyFollowsError{username, follower}
	}

	defer ub.resetCountCache(username, follower)
	_, err = conn.Exec(context.Background(), "INSERT INTO followers (username, followed_by) VALUES ($1, $2)",
		username, follower)
	if err != nil {
		return fmt.Errorf("could not insert into followers: %w", err)
	}

	return nil
}

// UnfollowUser ensures that username is unfollowed by follower,
// i.e., follower unfollows username
func (ub *userBacked) UnfollowUser(username, follower string) error {
	if username == follower {
		return fmt.Errorf("%s cannot unfollow itself", username)
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	isFollowingUser, err := ub.IsFollowingUser(username, follower)
	if err != nil {
		return fmt.Errorf("could not check follower: %w", err)
	}

	if !isFollowingUser {
		return UserDoesNotFollowError{username, follower}
	}

	defer ub.resetCountCache(username, follower)
	_, err = conn.Exec(context.Background(), "DELETE FROM followers WHERE username = $1 AND followed_by = $2",
		username, follower)
	if err != nil {
		return fmt.Errorf("could not delete row from followers: %w", err)
	}

	return nil
}

// IsFollowingUser checks if username is followed by follower,
// i.e., follower follows username
func (ub *userBacked) IsFollowingUser(username, follower string) (bool, error) {
	// TODO: this could be cached as well
	conn, err := ub.db.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), isFollowerOf, username, follower)
	if err != nil {
		return false, fmt.Errorf("could not perform isFollowerOf query: %w", err)
	}

	var countRows int
	for rows.Next() {
		if err = rows.Scan(&countRows); err != nil {
			return false, fmt.Errorf("could not scan isFollowerOf rows: %w", err)
		}
	}

	return countRows == 1, nil
}

// getUserDetails returns a PosterrUser containing
// the username and the date they joined
func (ub *userBacked) getUserDetails(username string) (types.PosterrUser, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return types.PosterrUser{}, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), selectUser, username)
	if err != nil {
		return types.PosterrUser{}, fmt.Errorf("could not perform selectUser query: %w", err)
	}

	var userProfile types.PosterrUser
	for rows.Next() {
		if err = rows.Scan(&userProfile.Username, &userProfile.JoinedAt); err != nil {
			return types.PosterrUser{}, fmt.Errorf("could not scan selectUser rows: %w", err)
		}
	}

	return userProfile, nil
}

// resetCountCache resets the counter cache for
// userA following count and userB followers count
func (ub *userBacked) resetCountCache(userA, userB string) {
	ub.Lock()
	defer ub.Unlock()
	// reset userB followers cache
	delete(ub.followersCount, userB)
	// reset userA following cache
	delete(ub.followingCount, userA)
}
