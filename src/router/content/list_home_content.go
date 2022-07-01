package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	dto := HomePageContentDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	posts, err := h.posts.ListHomePageContent(dto.Username, offset, dto.Toggle)
	if err != nil {
		switch err.(type) {
		case storage.InvalidToggleError:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete follow/unfollow operation: %s", err.Error())
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
