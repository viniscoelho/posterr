package content

import (
	"net/http"
	"strconv"

	"posterr/src/storage"
)

func parseIntQueryParam(param string) (int, error) {
	if len(param) != 0 {
		return strconv.Atoi(param)
	}
	return 0, nil
}

func getStatusCodeFromError(err error) int {
	switch err.(type) {
	case storage.PostExceededMaximumCharsError, storage.InvalidToggleError:
		return http.StatusBadRequest
	case storage.UserDoesNotExistError, storage.PostIdDoesNotExistError:
		return http.StatusNotFound
	case storage.ExceededMaximumDailyPostsError:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
