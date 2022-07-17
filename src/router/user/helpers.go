package user

import (
	"net/http"

	storageusers "posterr/src/storage/users"
)

const (
	targetUsernameQuery = "target"
)

func parseQueryParam(param string, r *http.Request) string {
	paramValue, exists := r.Form[param]
	if exists {
		return paramValue[0]
	}
	return ""
}

func getStatusCodeFromError(err error) int {
	switch err.(type) {
	case storageusers.SelfFollowError,
		storageusers.UserAlreadyFollowsError, storageusers.UserDoesNotFollowError:
		return http.StatusBadRequest
	case storageusers.UserDoesNotExistError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
