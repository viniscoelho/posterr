package content

type PostContentDTO struct {
	Username   string `json:"username"`
	Content    string `json:"content"`
	RepostedID string `json:"reposted_id"`
}
