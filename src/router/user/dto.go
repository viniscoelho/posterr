package user

type FollowToggle int

const (
	Follow FollowToggle = iota
	Unfollow
)

type FollowDTO struct {
	Username string       `json:"username"`
	Toggle   FollowToggle `json:"toggle"`
}
