package content

import (
	"encoding/json"
	"fmt"
	"net/http"

	"posterr/src/types"

	"github.com/sirupsen/logrus"
)

type searchContent struct {
	posts  types.Posterr
	logger *logrus.Entry
}

func NewSearchContentHandler(posts types.Posterr) *searchContent {
	return &searchContent{
		posts:  posts,
		logger: logrus.WithFields(logrus.Fields{"routes": "SearchContent"}),
	}
}

func (h *searchContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Errorf("Error parsing form: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))
	}

	limit, err := parseIntQueryParam(limitQuery, r)
	if err != nil {
		h.logger.Errorf("Error parsing limit: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	offset, err := parseIntQueryParam(offsetQuery, r)
	if err != nil {
		h.logger.Errorf("Error parsing offset: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	text := parseQueryParam(textQuery, r)

	posts, err := h.posts.SearchContent(text, limit, offset)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete search content operation: %s", err.Error())
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
