package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"posterr/src/storage"
	"posterr/src/types"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type readUser struct {
	users  types.Users
	logger *logrus.Entry
}

func NewReadUserHandler(users types.Users) *readUser {
	return &readUser{
		users:  users,
		logger: logrus.WithFields(logrus.Fields{"routes": "ReadUser"}),
	}
}

func (h *readUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := h.users.GetUserProfile(username)
	if err != nil {
		switch err.(type) {
		case storage.UserDoesNotExistError:
			rw.WriteHeader(http.StatusNotFound)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not get user details: %s", err)
		rw.Write([]byte(message))

		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	rw.Write(userBytes)
}
