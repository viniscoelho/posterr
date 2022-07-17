package posterr

import "fmt"

const (
	valueTooLongErrorCode        = "SQLSTATE 22001"
	foreignKeyViolationErrorCode = "SQLSTATE 23503"
)

type UserDoesNotExistError struct {
	username string
}

func (e UserDoesNotExistError) Error() string {
	return fmt.Sprintf("username %s is not registered", e.username)
}

type PostIdDoesNotExistError struct {
	postId string
}

func (e PostIdDoesNotExistError) Error() string {
	return fmt.Sprintf("post id %s is not registered", e.postId)
}

type ExceededMaximumDailyPostsError struct{}

func (e ExceededMaximumDailyPostsError) Error() string {
	return "exceeded maximum daily posts"
}

type PostExceededMaximumCharsError struct{}

func (e PostExceededMaximumCharsError) Error() string {
	return "post exceeded maximum allowed chars"
}

type InvalidToggleError struct{}

func (e InvalidToggleError) Error() string {
	return "invalid toggle selected"
}
