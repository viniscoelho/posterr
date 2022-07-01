package content

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"posterr/src/storage"
	"posterr/src/types"
)

type createContent struct {
	posts types.Posterr
}

func NewCreateContentHandler(posts types.Posterr) *createContent {
	return &createContent{posts}
}

func (h *createContent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	dto := PostContentDTO{}
	err = json.Unmarshal(body, &dto)
	if err != nil {
		log.Printf("CreateContentHandler request failed: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("internal server error"))

		return
	}

	_, err = h.posts.WriteContent(dto.Username, dto.Content, dto.RepostedID)
	if err != nil {
		switch err.(type) {
		case storage.ExceededMaximumDailyPostsError:
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("CreateContentHandler request failed: %s", err)
		message := fmt.Sprintf("could not complete write content operation: %s", err.Error())
		rw.Write([]byte(message))

		return
	}

	rw.WriteHeader(http.StatusCreated)
}
