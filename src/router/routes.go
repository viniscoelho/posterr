package router

import (
	"net/http"

	routercontent "posterr/src/router/content"
	routeruser "posterr/src/router/user"
	"posterr/src/types"

	"github.com/gorilla/mux"
)

func CreateRoutes(posts types.Posterr, users types.Users) *mux.Router {
	r := mux.NewRouter()

	r.Path("/posterr/content").
		Methods(http.MethodPost).
		Name("CreatePost").
		Handler(routercontent.NewCreateContentHandler(posts))
	r.Path("/posterr/content/home").
		Methods(http.MethodGet).
		Name("ListHomePosts").
		Handler(routercontent.NewListHomeContentHandler(posts))
	r.Path("/posterr/content/{username}").
		Methods(http.MethodGet).
		Name("ListProfileContent").
		Handler(routercontent.NewListProfileContentHandler(posts))

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
