package users

import (
	"strings"
)

func getErrorFromString(err error, username string) error {
	if strings.Contains(err.Error(), valueTooLongErrorCode) {
		return UsernameExceededMaximumCharsError{username}
	} else if strings.Contains(err.Error(), duplicatedKeyErrorCode) {
		return UserAlreadyExistsError{}
	} else if strings.Contains(err.Error(), noRowsInResult) {
		return UserDoesNotExistError{username}
	}
	return err
}
