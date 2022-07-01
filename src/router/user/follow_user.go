package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"posterr/src/storage"
	"posterr/src/types"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type followUser struct {
	users  types.Users
	logger *logrus.Entry
}

func NewFollowUserHandler(users types.Users) *followUser {
	return &followUser{
		users:  users,
		logger: logrus.WithFields(logrus.Fields{"routes": "FollowUser"}),
	}
}

func (h *followUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	dto := FollowDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	switch dto.Toggle {
	case Follow:
		err = h.users.FollowUser(username, dto.Username)
	case Unfollow:
		err = h.users.UnfollowUser(username, dto.Username)
	default:
		err = InvalidToggleError{}
	}
	if err != nil {
		switch err.(type) {
		case storage.SelfFollowError, storage.UserAlreadyFollowsError, storage.UserDoesNotFollowError:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete follow/unfollow operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
