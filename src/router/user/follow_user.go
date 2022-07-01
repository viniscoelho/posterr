package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"posterr/src/storage"
	"posterr/src/types"

	"github.com/gorilla/mux"
)

type followUser struct {
	users types.Users
}

func NewFollowUserHandler(users types.Users) *followUser {
	return &followUser{users}
}

func (h *followUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("FollowUserHandler failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	dto := FollowDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Printf("FollowUserHandler failed: %s", err)
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
	}
	if err != nil {
		switch err.(type) {
		case storage.UserAlreadyFollowsError, storage.UserDoesNotFollowError:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("FollowUserHandler failed: %s", err)
		message := fmt.Sprintf("could not complete follow/unfollow operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
