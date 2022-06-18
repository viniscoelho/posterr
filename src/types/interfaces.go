//go:generate mockgen -destination=mocks/mocks.go -package=mocks posterr/src/types Posterr,Users
package types

import "time"

type PostsListToggle int

const (
	All PostsListToggle = iota
	Following
)

type PosterrUser struct {
	Username  string    `json:"username"`
	Followers int       `json:"followers"`
	Following int       `json:"following"`
	Posts     int       `json:"posts"`
	JoinedAt  time.Time `json:"joined_at"`
}

type PosterrContent struct {
	ID         int       `json:"post_id"`
	Username   string    `json:"username"`
	Content    string    `json:"content,omitempty"`
	RepostedId int       `json:"reposted_id,omitempty"`
	CreatedOn  time.Time `json:"created_on"`
}

type Posterr interface {
	ListHomePagePosts(username string, offset int, toggle PostsListToggle) ([]PosterrContent, error)
	ListProfilePosts(username string, offset int) ([]PosterrContent, error)
	WritePost(username, postContent string, repostedId int) error
}

type User interface {
	GetUserProfile(username string) (PosterrUser, error)
	CountUserPosts(username string) (int, error)
	CountUserFollowers(username string) (int, error)
	CountUserFollowing(username string) (int, error)
	FollowUser(followerUsername, followingUsername string) error
	UnfollowUser(followerUsername, followingUsername string) error
	IsFollowingUser(followerUsername, followingUsername string) (bool, error)
}
