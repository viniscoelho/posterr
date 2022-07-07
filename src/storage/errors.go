package storage

import "fmt"

const (
	valueTooLongErrorCode        = "SQLSTATE 22001"
	foreignKeyViolationErrorCode = "SQLSTATE 23503"
	noRowsInResult               = "no rows in result set"
)

type UserDoesNotExistError struct {
	username string
}

func (e UserDoesNotExistError) Error() string {
	return fmt.Sprintf("username %s is not registered", e.username)
}

type InvalidUsernameError struct {
	username string
}

func (e InvalidUsernameError) Error() string {
	return fmt.Sprintf("invalid username %s: username must consist of alphanumeric charactes only", e.username)
}

type PostIdDoesNotExistError struct {
	postId string
}

func (e PostIdDoesNotExistError) Error() string {
	return fmt.Sprintf("post id %s is not registered", e.postId)
}

type SelfFollowError struct {
	user string
}

func (e SelfFollowError) Error() string {
	return fmt.Sprintf("%s cannot unfollow itself", e.user)
}

type UserAlreadyFollowsError struct {
	user     string
	follower string
}

func (e UserAlreadyFollowsError) Error() string {
	return fmt.Sprintf("%s already follows %s", e.follower, e.user)
}

type UserDoesNotFollowError struct {
	user     string
	follower string
}

func (e UserDoesNotFollowError) Error() string {
	return fmt.Sprintf("%s does not follow %s", e.follower, e.user)
}

type ExceededMaximumDailyPostsError struct{}

func (e ExceededMaximumDailyPostsError) Error() string {
	return "exceeded maximum daily posts"
}

type InvalidToggleError struct{}

func (e InvalidToggleError) Error() string {
	return "invalid toggle selected"
}

type UsernameExceededMaximumCharsError struct {
	username string
}

func (e UsernameExceededMaximumCharsError) Error() string {
	return fmt.Sprintf("username %s exceeded maximum allowed chars", e.username)
}

type PostExceededMaximumCharsError struct{}

func (e PostExceededMaximumCharsError) Error() string {
	return "post exceeded maximum allowed chars"
}
