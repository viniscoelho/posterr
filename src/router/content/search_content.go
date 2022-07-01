package content

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	username := r.FormValue("username")
	limitQuery := r.FormValue("limit")
	offsetQuery := r.FormValue("offset")

	limit, err := h.parseIntQueryParams(limitQuery)
	if err != nil {
		h.logger.Errorf("Error parsing limit: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	offset, err := h.parseIntQueryParams(offsetQuery)
	if err != nil {
		h.logger.Errorf("Error parsing offset: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	posts, err := h.posts.SearchContent(username, limit, offset)
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

func (h searchContent) parseIntQueryParams(paramStr string) (int, error) {
	if len(paramStr) != 0 {
		return strconv.Atoi(paramStr)
	}
	return 0, nil
}
