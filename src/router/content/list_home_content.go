package content

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	username := parseQueryParam(usernameQuery, r)
	toggle := parseBoolQueryParam(toggleQuery, r)

	posts, err := h.posts.ListHomePageContent(username, offset, toggle)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		rw.WriteHeader(statusCode)
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
