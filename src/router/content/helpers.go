package content

import (
	"net/http"
	"strconv"

	storageposterr "posterr/src/storage/posterr"
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
	case storageposterr.PostExceededMaximumCharsError, storageposterr.InvalidToggleError:
		return http.StatusBadRequest
	case storageposterr.UserDoesNotExistError, storageposterr.PostIdDoesNotExistError:
		return http.StatusNotFound
	case storageposterr.ExceededMaximumDailyPostsError:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
