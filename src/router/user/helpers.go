package user

import (
	"net/http"

	"posterr/src/storage"
)

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
