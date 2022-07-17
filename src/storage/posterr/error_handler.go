package posterr

import (
	"strings"
)

func getErrorFromString(err error, username, repostedId string) error {
	if strings.Contains(err.Error(), valueTooLongErrorCode) {
		return PostExceededMaximumCharsError{}
	} else if strings.Contains(err.Error(), foreignKeyViolationErrorCode) {
		if strings.Contains(err.Error(), "posts_username_fkey") {
			return UserDoesNotExistError{username}
		} else if strings.Contains(err.Error(), "posts_reposted_id_fkey") {
			return PostIdDoesNotExistError{repostedId}
		}
	}
	return err
}
