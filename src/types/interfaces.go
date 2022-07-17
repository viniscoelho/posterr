//go:generate mockgen -destination=mocks/mocks.go -package=mocks posterr/src/types Posterr,Users
package types

import (
	"time"
)

const (
	All       = false
	Following = true
)

type PosterrUser struct {
	Username string `json:"username"`
}

type PosterrUserDetailed struct {
	PosterrUser
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
	ListHomePageContent(username string, offset int, toggle bool) ([]PosterrContent, error)
	ListProfileContent(username string, offset int) ([]PosterrContent, error)
	SearchContent(text string, limit, offset int) ([]PosterrContent, error)
	WriteContent(username, postContent string) (string, error)
	WriteRepostContent(username, repostedId string) (string, error)
	WriteQuoteRepostContent(username, postContent, repostedId string) (string, error)
}

type Users interface {
	CreateUser(username string) error
	GetUserProfile(username string) (PosterrUserDetailed, error)
	CountUserPosts(username string) (int, error)
	CountUserFollowers(username string) (int, error)
	CountUserFollowing(username string) (int, error)
	ListFollowers(username string) ([]PosterrUser, error)
	FollowUser(targetUser, currentUser string) error
	UnfollowUser(targetUser, currentUser string) error
	IsFollowingUser(targetUser, currentUser string) (bool, error)
}
