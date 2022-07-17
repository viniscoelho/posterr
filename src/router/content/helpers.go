package content

import (
	"net/http"
	"strconv"

	"posterr/src/storage"
)

const (
	usernameQuery = "username"
	limitQuery    = "limit"
	offsetQuery   = "offset"
	textQuery     = "text"
	toggleQuery   = "toggle"
)

func parseQueryParam(param string, r *http.Request) string {
	paramValue, exists := r.Form[param]
	if exists {
		return paramValue[0]
	}
	return ""
}

func parseIntQueryParam(param string, r *http.Request) (int, error) {
	paramValue, exists := r.Form[param]
	if exists {
		value, err := strconv.Atoi(paramValue[0])
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, nil
}

func parseBoolQueryParam(param string, r *http.Request) bool {
	_, exists := r.Form[param]
	return exists
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
