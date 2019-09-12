package main

import "net/http"

type Route struct {
	Name		string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// We use meaningful HTTP verbs for the API.
// This is needed because we don't want, for example,
// to use GET requests for something that change the
// server state (e.g. stateless).

var routes = []Route{

	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"UserAdd",
		"POST",
		"/api/users/",
		UserAdd,
	},

	Route{
		"UserEdit",
		"PUT",
		"/api/users/{user}",
		UserEdit,
	},

	Route{
		"UserDel",
		"DELETE",
		"/api/users/{user}",
		UserDel,
	},

	Route{
		"UserSearch",
		"GET",
		"/api/users/search",
		UserSearch,
	},

}
