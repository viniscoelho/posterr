package storage

import (
	"context"
	"fmt"
	"posterr/src/storage/postgres"
	"posterr/src/types"

	"github.com/jackc/pgx/v4"
)

const (
	maxDailyPosts = 5

	selectAllPosts = `SELECT *
                 FROM posts
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $1`

	selectFollowingPosts = `SELECT *
                 FROM posts
                 WHERE username IN (
                     SELECT followed_by
                     FROM follows 
                     WHERE followed_by = $1) 
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $2`

	selectProfilePosts = `SELECT *
                 FROM posts
                 WHERE username = $1
                 ORDER BY created_at DESC
                 LIMIT 5
                 OFFSET $2`

	countDailyPosts = `SELECT COUNT(*) as daily_posts
                 FROM posts
                 WHERE username = $1
                 AND date_trunc('day', created_at) = date_trunc('day', NOW())`
)

type posterrBacked struct {
	db postgres.ConnectDB
}

func NewPosterrBacked(db postgres.ConnectDB) *posterrBacked {
	return &posterrBacked{
		db: db,
	}
}

func (pb *posterrBacked) ListHomePagePosts(username string, offset int, toggle types.PostsListToggle) ([]types.PosterrContent, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	var rows pgx.Rows
	switch toggle {
	case types.All:
		rows, err = conn.Query(context.Background(), selectAllPosts, offset)
		if err != nil {
			return nil, fmt.Errorf("could not perform query: %w", err)
		}
	case types.Following:
		rows, err = conn.Query(context.Background(), selectFollowingPosts, username, offset)
		if err != nil {
			return nil, fmt.Errorf("could not perform query: %w", err)
		}
	}

	posts := make([]types.PosterrContent, 0)
	for rows.Next() {
		postContent := types.PosterrContent{}
		if err = rows.Scan(&postContent.ID, &postContent.Username, &postContent.Content, &postContent.RepostedId, &postContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan rows: %w", err)
		}

		posts = append(posts, postContent)
	}

	return posts, nil
}

func (pb *posterrBacked) ListProfilePosts(username string, offset int) ([]types.PosterrContent, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), selectProfilePosts, username, offset)
	if err != nil {
		return nil, fmt.Errorf("could not perform query: %w", err)
	}

	posts := make([]types.PosterrContent, 0)
	for rows.Next() {
		postContent := types.PosterrContent{}
		if err = rows.Scan(&postContent.ID, &postContent.Username, &postContent.Content, &postContent.RepostedId, &postContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan rows: %w", err)
		}

		posts = append(posts, postContent)
	}

	return posts, nil
}

func (pb *posterrBacked) WritePost(username, postContent string, repostedId int) error {
	// TODO: Content size should not be superior than 777 chars
	if len(postContent) == 0 && repostedId == 0 {
		return fmt.Errorf("either content or reposted_id should have a value")
	}

	conn, err := pb.db.Connect()
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	dailyPosts, err := pb.countDailyPosts(username)
	if err != nil {
		return fmt.Errorf("could not count daily posts: %w", err)
	}

	if dailyPosts >= maxDailyPosts {
		return fmt.Errorf("exceeded maximum daily posts")
	} else if repostedId == 0 {
		// if repostedId is zero, this is a regular post
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (username, content) VALUES ($1, $2)",
			username, postContent)
	} else if len(postContent) == 0 {
		// if postContent is empty, this is a repost
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (username, reposted_id) VALUES ($1, $2)",
			username, repostedId)
	} else {
		// otherwise, this is a quoted-repost
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (username, content, reposted_id) VALUES ($1, $2, $3)",
			username, postContent, repostedId)
	}

	if err != nil {
		return fmt.Errorf("could not insert into posts: %w", err)
	}

	return nil
}

func (pb *posterrBacked) countDailyPosts(username string) (int, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countDailyPosts, username)
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
