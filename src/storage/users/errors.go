package users

import "fmt"

const (
	valueTooLongErrorCode  = "SQLSTATE 22001"
	duplicatedKeyErrorCode = "SQLSTATE 23505"
	noRowsInResult         = "no rows in result set"
)

type InvalidUsernameError struct {
	username string
}

func (e InvalidUsernameError) Error() string {
	return fmt.Sprintf("invalid username %s: username must consist of alphanumeric charactes only", e.username)
}

type UsernameExceededMaximumCharsError struct {
	username string
}

func (e UsernameExceededMaximumCharsError) Error() string {
	return fmt.Sprintf("username %s exceeded maximum allowed chars", e.username)
}

type UserAlreadyExistsError struct {
	username string
}

func (e UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("username %s already exists", e.username)
}

type UserDoesNotExistError struct {
	username string
}

func (e UserDoesNotExistError) Error() string {
	return fmt.Sprintf("username %s is not registered", e.username)
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
