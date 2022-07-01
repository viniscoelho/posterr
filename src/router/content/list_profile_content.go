package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"posterr/src/types"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type listProfileContent struct {
	posts  types.Posterr
	logger *logrus.Entry
}

func NewListProfileContentHandler(posts types.Posterr) *listProfileContent {
	return &listProfileContent{
		posts:  posts,
		logger: logrus.WithFields(logrus.Fields{"routes": "ListProfileContent"}),
	}
}

func (h *listProfileContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
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

	profilePosts, err := h.posts.ListProfileContent(username, offset)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete list profile content operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	postsBytes, err := json.Marshal(profilePosts)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	rw.Write(postsBytes)
}
