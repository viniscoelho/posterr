package user

import (
	"encoding/json"
	"net/http"
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
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

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
