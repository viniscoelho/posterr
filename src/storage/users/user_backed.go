package users

import (
	"context"
	"fmt"
	"regexp"
	"sync"

	storagedb "posterr/src/storage/db"
	"posterr/src/types"
)

type userBacked struct {
	sync.RWMutex
	// An accessor to the database
	db storagedb.ConnectDB
	// A cache map to store how many followers a user has
	followersCount map[string]int
	// A cache map to store how many users a user is following
	followingCount map[string]int
	// A regex to validate usernames
	rgx *regexp.Regexp
}

func NewUserBacked(db storagedb.ConnectDB) *userBacked {
	// TODO: should the connection pool be reset for every call?
	return &userBacked{
		db:             db,
		followersCount: make(map[string]int),
		followingCount: make(map[string]int),
		rgx:            regexp.MustCompile(`^[a-zA-Z0-9]*$`),
	}
}

// CreateUser creates a user. This method is not exposed through an API
func (ub *userBacked) CreateUser(username string) error {
	if !ub.rgx.MatchString(username) {
		return InvalidUsernameError{username}
	}

	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), "INSERT INTO users (username) VALUES ($1)", username)
	if err != nil {
		err = fmt.Errorf("could not insert into users: %w", err)
		return getErrorFromString(err, username)
	}

	return nil
}

// GetUserProfile returns a user detailed information
func (ub *userBacked) GetUserProfile(username string) (types.PosterrUserDetailed, error) {
	userProfile, err := ub.getUserDetails(username)
	if err != nil {
		return types.PosterrUserDetailed{}, err
	}

	userProfile.Followers, err = ub.CountUserFollowers(username)
	if err != nil {
		return types.PosterrUserDetailed{}, err
	}

	userProfile.Following, err = ub.CountUserFollowing(username)
	if err != nil {
		return types.PosterrUserDetailed{}, err
	}

	userProfile.PostsCount, err = ub.CountUserPosts(username)
	if err != nil {
		return types.PosterrUserDetailed{}, err
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

	var dailyPosts int
	row := conn.QueryRow(context.Background(), countUserPosts, username)
	if err = row.Scan(&dailyPosts); err != nil {
		return 0, fmt.Errorf("could not scan countUserPosts rows: %w", err)
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

	var followers int
	row := conn.QueryRow(context.Background(), countFollowers, username)
	if err = row.Scan(&followers); err != nil {
		return 0, fmt.Errorf("could not scan countFollowers rows: %w", err)
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

	var following int
	row := conn.QueryRow(context.Background(), countFollowing, username)
	if err = row.Scan(&following); err != nil {
		return 0, fmt.Errorf("could not scan countFollowing rows: %w", err)
	}

	func(username string) {
		ub.Lock()
		defer ub.Unlock()
		ub.followingCount[username] = following
	}(username)

	return following, nil
}

// ListFollowers returns a list of followers of a user
func (ub *userBacked) ListFollowers(username string) ([]types.PosterrUser, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	_, err = ub.getUserDetails(username)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(context.Background(), listFollowers, username)
	if err != nil {
		return nil, fmt.Errorf("could not perform listFollowers query: %w", err)
	}

	followers := make([]types.PosterrUser, 0)
	for rows.Next() {
		var follower types.PosterrUser
		if err = rows.Scan(&follower.Username); err != nil {
			return nil, fmt.Errorf("could not scan listFollowers rows: %w", err)
		}

		followers = append(followers, follower)
	}

	return followers, nil
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
		return SelfFollowError{username}
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

	var countRows int
	row := conn.QueryRow(context.Background(), isFollowerOf, username, follower)
	if err = row.Scan(&countRows); err != nil {
		return false, fmt.Errorf("could not scan isFollowerOf rows: %w", err)
	}

	return countRows == 1, nil
}

// getUserDetails returns a PosterrUser containing
// the username and the date they joined
func (ub *userBacked) getUserDetails(username string) (types.PosterrUserDetailed, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return types.PosterrUserDetailed{}, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	var userProfile types.PosterrUserDetailed
	row := conn.QueryRow(context.Background(), selectUser, username)
	if err = row.Scan(&userProfile.Username, &userProfile.JoinedAt); err != nil {
		err = fmt.Errorf("could not scan selectUser rows: %w", err)
		return types.PosterrUserDetailed{}, getErrorFromString(err, username)
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
