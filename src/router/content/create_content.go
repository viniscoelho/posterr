package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"posterr/src/types"

	"github.com/sirupsen/logrus"
)

type createContent struct {
	posts  types.Posterr
	logger *logrus.Entry
}

func NewCreateContentHandler(posts types.Posterr) *createContent {
	return &createContent{
		posts:  posts,
		logger: logrus.WithFields(logrus.Fields{"routes": "CreateContent"}),
	}
}

func (h *createContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	dto := PostContentDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		h.logger.Errorf("Request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	_, err = h.posts.WriteContent(dto.Username, dto.Content, dto.RepostedID)
	if err != nil {
		statusCode := getStatusCodeFromError(err)
		rw.WriteHeader(statusCode)
		h.logger.Errorf("Request failed: %s", err)
		message := fmt.Sprintf("could not complete write content operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	rw.WriteHeader(http.StatusCreated)
}
