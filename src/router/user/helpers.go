package user

import (
	"net/http"
	"posterr/src/storage"
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
	case storage.SelfFollowError,
		storage.UserAlreadyFollowsError, storage.UserDoesNotFollowError:
		return http.StatusBadRequest
	case storage.UserDoesNotExistError, storage.PostIdDoesNotExistError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
