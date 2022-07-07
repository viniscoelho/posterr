package content

import (
	"encoding/json"
	"fmt"
	"net/http"

	"posterr/src/storage"
	"posterr/src/types"

	"github.com/sirupsen/logrus"
)

type listHomeContent struct {
	posts  types.Posterr
	logger *logrus.Entry
}

func NewListHomeContentHandler(posts types.Posterr) *listHomeContent {
	return &listHomeContent{
		posts:  posts,
		logger: logrus.WithFields(logrus.Fields{"routes": "ListHomeContent"}),
	}
}

func (h *listHomeContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	offsetQuery := r.FormValue("offset")
	toggleQuery := r.FormValue("toggle")

	offset, err := parseIntQueryParam(offsetQuery)
	if err != nil {
		h.logger.Errorf("Error parsing limit: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}

	toggle, err := parseIntQueryParam(toggleQuery)
	if err != nil {
		h.logger.Errorf("Error parsing toggle: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}

	posts, err := h.posts.ListHomePageContent(username, offset, types.PostsListToggle(toggle))
	if err != nil {
		switch err.(type) {
		case storage.InvalidToggleError:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete list home page content operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	postsBytes, err := json.Marshal(posts)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	rw.Write(postsBytes)
}
