package user

import (
	"fmt"
	"net/http"

	"posterr/src/types"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type unfollowUser struct {
	users  types.Users
	logger *logrus.Entry
}

func NewUnfollowUserHandler(users types.Users) *unfollowUser {
	return &unfollowUser{
		users:  users,
		logger: logrus.WithFields(logrus.Fields{"routes": "UnfollowUser"}),
	}
}

func (h *unfollowUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("Error parsing form: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}

	vars := mux.Vars(r)
	username := vars["username"]
	targetUsername := parseQueryParam(targetUsernameQuery, r)

	err = h.users.UnfollowUser(username, targetUsername)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		rw.WriteHeader(statusCode)
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete follow/unfollow operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
