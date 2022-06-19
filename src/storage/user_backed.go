package storage

import (
	"context"
	"fmt"
	"posterr/src/postgres"
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

type userBacked struct {
	db    postgres.ConnectDB
	posts types.Posterr
}

func NewUserBacked(db postgres.ConnectDB, posts types.Posterr) *userBacked {
	return &userBacked{
		db:    db,
		posts: posts,
	}
}

func (ub *userBacked) CreateUser(username string) error {
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

	userProfile.Posts, err = ub.posts.ListProfilePosts(username, offset)
	if err != nil {
		return types.PosterrUser{}, err
	}

	return userProfile, nil
}

func (ub *userBacked) CountUserPosts(username string) (int, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countUserPosts, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform query: %w", err)
	}

	var dailyPosts int
	for rows.Next() {
		if err = rows.Scan(&dailyPosts); err != nil {
			return 0, fmt.Errorf("could not scan rows: %w", err)
		}
	}

	return dailyPosts, nil
}

func (ub *userBacked) CountUserFollowers(username string) (int, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countFollowers, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform query: %w", err)
	}

	var followers int
	for rows.Next() {
		if err = rows.Scan(&followers); err != nil {
			return 0, fmt.Errorf("could not scan rows: %w", err)
		}
	}

	return followers, nil
}

func (ub *userBacked) CountUserFollowing(username string) (int, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countFollowing, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform query: %w", err)
	}

	var following int
	for rows.Next() {
		if err = rows.Scan(&following); err != nil {
			return 0, fmt.Errorf("could not scan rows: %w", err)
		}
	}

	return following, nil
}

func (ub *userBacked) FollowUser(followerUsername, followingUsername string) error {
	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	isFollowingUser, err := ub.IsFollowingUser(followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("could not check activity: %w", err)
	}

	if isFollowingUser {
		return fmt.Errorf("%s already follows %s", followerUsername, followingUsername)
	}

	_, err = conn.Exec(context.Background(), "INSERT INTO followers (username, followed_by) VALUES ($1, $2)",
		followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("could not insert into followers: %w", err)
	}

	return nil
}

func (ub *userBacked) UnfollowUser(followerUsername, followingUsername string) error {
	conn, err := ub.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	isFollowingUser, err := ub.IsFollowingUser(followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("could not check activity: %w", err)
	}

	if !isFollowingUser {
		return fmt.Errorf("%s does not follow %s", followerUsername, followingUsername)
	}

	_, err = conn.Exec(context.Background(), "DELETE FROM followers WHERE username = $1 AND followed_by = $2",
		followerUsername, followingUsername)
	if err != nil {
		return fmt.Errorf("could not delete row from followers: %w", err)
	}

	return nil
}

func (ub *userBacked) IsFollowingUser(followerUsername, followingUsername string) (bool, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return false, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), isFollowerOf, followerUsername, followingUsername)
	if err != nil {
		return false, fmt.Errorf("could not perform query: %w", err)
	}

	var countRows int
	for rows.Next() {
		if err = rows.Scan(&countRows); err != nil {
			return false, fmt.Errorf("could not scan rows: %w", err)
		}
	}

	return countRows == 1, nil
}

func (ub *userBacked) getUserDetails(username string) (types.PosterrUser, error) {
	conn, err := ub.db.Connect()
	if err != nil {
		return types.PosterrUser{}, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), selectUser, username)
	if err != nil {
		return types.PosterrUser{}, fmt.Errorf("could not perform query: %w", err)
	}

	var userProfile types.PosterrUser
	for rows.Next() {
		if err = rows.Scan(&userProfile.Username, &userProfile.JoinedAt); err != nil {
			return types.PosterrUser{}, fmt.Errorf("could not scan rows: %w", err)
		}
	}

	return userProfile, nil
}
