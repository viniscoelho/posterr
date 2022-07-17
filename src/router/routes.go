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
		Name("CreateContent").
		Handler(routercontent.NewCreateContentHandler(posts))
	r.Path("/posterr/content").
		Methods(http.MethodGet).
		Name("SearchContent").
		Handler(routercontent.NewSearchContentHandler(posts))
	r.Path("/posterr/content/home").
		Methods(http.MethodGet).
		Name("ListHomeContent").
		Handler(routercontent.NewListHomeContentHandler(posts))
	r.Path("/posterr/content/{username}").
		Methods(http.MethodGet).
		Name("ListProfileContent").
		Handler(routercontent.NewListProfileContentHandler(posts))

	r.Path("/posterr/users/{username}").
		Methods(http.MethodGet).
		Name("ReadUser").
		Handler(routeruser.NewReadUserHandler(users))
	r.Path("/posterr/users/{username}/followers").
		Methods(http.MethodGet).
		Name("ListFollowers").
		Handler(routeruser.NewListFollowersHandler(users))

	r.Path("/posterr/users/{username}/follow").
		Methods(http.MethodPost).
		Name("FollowUser").
		Handler(routeruser.NewFollowUserHandler(users))
	r.Path("/posterr/users/{username}/unfollow").
		Methods(http.MethodPost).
		Name("UnfollowUser").
		Handler(routeruser.NewUnfollowUserHandler(users))

	return r
}
