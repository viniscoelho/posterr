package content

import "posterr/src/types"

type PostContentDTO struct {
	Username   string `json:"username"`
	Content    string `json:"content"`
	RepostedID string `json:"reposted_id"`
}

type HomePageContentDTO struct {
	Username string                `json:"username"`
	Toggle   types.PostsListToggle `json:"toggle"`
}
