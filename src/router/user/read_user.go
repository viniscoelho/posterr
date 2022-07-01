package user

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	offsetStr := r.FormValue("offset")

	var err error
	var offset int
	if len(offsetStr) != 0 {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Errorf("Error parsing offset: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("internal server error"))

			return
		}
	}

	user, err := h.users.GetUserProfile(username, offset)
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
