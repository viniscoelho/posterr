package storage

import "fmt"

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
