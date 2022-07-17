package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"posterr/src/types"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type listFollowers struct {
	users  types.Users
	logger *logrus.Entry
}

func NewListFollowersHandler(users types.Users) *listFollowers {
	return &listFollowers{
		users:  users,
		logger: logrus.WithFields(logrus.Fields{"routes": "ListFollowers"}),
	}
}

func (h *listFollowers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	followers, err := h.users.ListFollowers(username)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		rw.WriteHeader(statusCode)
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete follow/unfollow operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	followersBytes, err := json.Marshal(followers)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	rw.Write(followersBytes)
}
