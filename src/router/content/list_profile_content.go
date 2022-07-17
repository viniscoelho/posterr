package content

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("Error parsing param form: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}

	offset, err := parseIntQueryParam(offsetQuery, r)
	if err != nil {
		h.logger.Errorf("Error parsing offset: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
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
