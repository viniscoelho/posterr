//go:generate mockgen -destination=mocks/mocks.go -package=mocks posterr/src/types Posterr,Users
package types

import (
	"time"
)

type PostsListToggle int

const (
	All PostsListToggle = iota
	Following
)

type PosterrUser struct {
	Username   string    `json:"username"`
	Followers  int       `json:"followers"`
	Following  int       `json:"following"`
	PostsCount int       `json:"posts_count"`
	JoinedAt   time.Time `json:"joined_at"`
}

type PosterrContent struct {
	ID         string    `json:"post_id,omitempty"`
	Username   string    `json:"username"`
	Content    string    `json:"content,omitempty"`
	RepostedId string    `json:"reposted_id,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

type Posterr interface {
	ListHomePageContent(username string, offset int, toggle PostsListToggle) ([]PosterrContent, error)
	ListProfileContent(username string, offset int) ([]PosterrContent, error)
	SearchContent(content string, limit, offset int) ([]PosterrContent, error)
	WriteContent(username, postContent, repostedId string) (string, error)
}

type Users interface {
	CreateUser(username string) error
	GetUserProfile(username string) (PosterrUser, error)
	CountUserPosts(username string) (int, error)
	CountUserFollowers(username string) (int, error)
	CountUserFollowing(username string) (int, error)
	FollowUser(userA, userB string) error
	UnfollowUser(userA, userB string) error
	IsFollowingUser(userA, userB string) (bool, error)
}
