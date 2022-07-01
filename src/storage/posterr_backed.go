package storage

import (
	"context"
	"fmt"

	storagedb "posterr/src/storage/db"
	"posterr/src/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

const (
	maxDailyPosts = 5

	selectAllPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $1`

	selectFollowingPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE username IN (
                     SELECT followed_by
                     FROM followers
                     WHERE followed_by = $1) 
                 ORDER BY created_at DESC
                 LIMIT 10
                 OFFSET $2`

	selectProfilePosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE username = $1
                 ORDER BY created_at DESC
                 LIMIT 5
                 OFFSET $2`

	countDailyPosts = `SELECT COUNT(*) as daily_posts
                 FROM posts
                 WHERE username = $1
                 AND date_trunc('day', created_at) = date_trunc('day', NOW())`

	searchPosts = `SELECT post_id, username, COALESCE(content, ''), COALESCE(reposted_id, ''), created_at
                 FROM posts
                 WHERE content LIKE '%$1%'
                 ORDER BY created_at DESC
                 LIMIT $2
                 OFFSET $3`
)

type posterrBacked struct {
	db storagedb.ConnectDB
}

func NewPosterrBacked(db storagedb.ConnectDB) *posterrBacked {
	return &posterrBacked{
		db: db,
	}
}

// ListHomePageContent returns a list of posts:
// - If the toggle is All, returns a list of posts from the whole database
// - If the toggle is Following, returns a list of posts only from the users a given username follows
// Each call returns 10 posts at most
func (pb *posterrBacked) ListHomePageContent(username string, offset int, toggle types.PostsListToggle) ([]types.PosterrContent, error) {
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
			return nil, fmt.Errorf("could not perform selectAllPosts query: %w", err)
		}
	case types.Following:
		rows, err = conn.Query(context.Background(), selectFollowingPosts, username, offset)
		if err != nil {
			return nil, fmt.Errorf("could not perform selectFollowingPosts query: %w", err)
		}
	default:
		return nil, InvalidToggleError{}
	}

	posts := make([]types.PosterrContent, 0)
	for rows.Next() {
		postContent := types.PosterrContent{}
		if err = rows.Scan(&postContent.ID, &postContent.Username, &postContent.Content, &postContent.RepostedId, &postContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan selectPosts rows: %w", err)
		}

		posts = append(posts, postContent)
	}

	return posts, nil
}

// ListProfileContent returns a lists of posts for a given username
// Each call returns 5 posts at most
func (pb *posterrBacked) ListProfileContent(username string, offset int) ([]types.PosterrContent, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), selectProfilePosts, username, offset)
	if err != nil {
		return nil, fmt.Errorf("could not perform selectProfilePosts query: %w", err)
	}

	posts := make([]types.PosterrContent, 0)
	for rows.Next() {
		postContent := types.PosterrContent{}
		if err = rows.Scan(&postContent.ID, &postContent.Username, &postContent.Content, &postContent.RepostedId, &postContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan selectProfilePosts rows: %w", err)
		}

		posts = append(posts, postContent)
	}

	return posts, nil
}

// SearchContent returns a lists of posts matching a substring criteria.
// The number of returned posts can be specified by the providing a limit.
func (pb *posterrBacked) SearchContent(content string, limit, offset int) ([]types.PosterrContent, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), searchPosts, content, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not perform searchPosts query: %w", err)
	}

	posts := make([]types.PosterrContent, 0)
	for rows.Next() {
		postContent := types.PosterrContent{}
		if err = rows.Scan(&postContent.ID, &postContent.Username, &postContent.Content, &postContent.RepostedId, &postContent.CreatedAt); err != nil {
			return nil, fmt.Errorf("could not scan searchPosts rows: %w", err)
		}

		posts = append(posts, postContent)
	}

	return posts, nil
}

// WriteContent creates a post for a given username and returns the postId
func (pb *posterrBacked) WriteContent(username, postContent, repostedId string) (string, error) {
	if len(postContent) == 0 && len(repostedId) == 0 {
		return "", fmt.Errorf("either content or reposted_id should have a value")
	}

	conn, err := pb.db.Connect()
	if err != nil {
		return "", fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	dailyPosts, err := pb.countDailyPosts(username)
	if err != nil {
		return "", fmt.Errorf("could not count daily posts: %w", err)
	}

	postId := uuid.New().String()
	if dailyPosts >= maxDailyPosts {
		return "", ExceededMaximumDailyPostsError{}
	} else if len(repostedId) == 0 {
		// if repostedId is empty, this is a regular post
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (post_id, username, content) VALUES ($1, $2, $3)",
			postId, username, postContent)
	} else if len(postContent) == 0 {
		// if postContent is empty, this is a repost
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (post_id, username, reposted_id) VALUES ($1, $2, $3)",
			postId, username, repostedId)
	} else {
		// otherwise, this is a quoted-repost
		_, err = conn.Exec(context.Background(), "INSERT INTO posts (post_id, username, content, reposted_id) VALUES ($1, $2, $3, $4)",
			postId, username, postContent, repostedId)
	}

	if err != nil {
		return "", fmt.Errorf("could not insert into posts: %w", err)
	}

	return postId, nil
}

// countDailyPosts returns how many posts where made in a single day
func (pb *posterrBacked) countDailyPosts(username string) (int, error) {
	conn, err := pb.db.Connect()
	if err != nil {
		return 0, fmt.Errorf("could not connect to database: %w", err)
	}
	defer conn.Close()

	rows, err := conn.Query(context.Background(), countDailyPosts, username)
	if err != nil {
		return 0, fmt.Errorf("could not perform countDailyPosts query: %w", err)
	}

	var dailyPosts int
	for rows.Next() {
		if err = rows.Scan(&dailyPosts); err != nil {
			return 0, fmt.Errorf("could not scan countDailyPosts rows: %w", err)
		}
	}

	return dailyPosts, nil
}
