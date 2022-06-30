package router

import (
	"net/http"

	routerposts "posterr/src/router/posts"
	routeruser "posterr/src/router/user"
	"posterr/src/types"

	"github.com/gorilla/mux"
)

func CreateRoutes(posts types.Posterr, users types.Users) *mux.Router {
	r := mux.NewRouter()

	r.Path("/posterr/content").
		Methods(http.MethodPost).
		Name("CreatePost").
		Handler(routerposts.NewCreatePostHandler(posts))
	r.Path("/posterr/content/home").
		Methods(http.MethodGet).
		Name("ListHomePosts").
		Handler(routerposts.NewListHomePostsHandler(posts))
	r.Path("/posterr/content/profile/{username}").
		Methods(http.MethodGet).
		Name("ListProfilePosts").
		Handler(routerposts.NewListProfilePostsHandler(posts))

	r.Path("/posterr/user/{username}").
		Methods(http.MethodGet).
		Name("ReadUser").
		Handler(routeruser.NewReadUserHandler(users))
	r.Path("/posterr/user/{username}/follow").
		Methods(http.MethodPost).
		Name("FollowUser").
		Handler(routeruser.NewFollowUserHandler(users))

	return r
}
