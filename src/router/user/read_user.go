package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"posterr/src/types"

	"github.com/gorilla/mux"
)

type readUser struct {
	users types.Users
}

func NewReadUserHandler(users types.Users) *readUser {
	return &readUser{users}
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
			log.Printf("Error parsing offset: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("internal server error"))

			return
		}
	}

	user, err := h.users.GetUserProfile(username, offset)
	if err != nil {
		log.Printf("ReadUserHandler request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		log.Printf("ReadUserHandler request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	rw.Write(userBytes)
}
