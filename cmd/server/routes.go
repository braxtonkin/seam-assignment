package main

import (
	"net/http"
)

// Initializes and returns a new http.ServeMux, with routes registered
func (app *application) routes() http.Handler {
	// Using  golang's 1.22 http routing
	mux := http.NewServeMux()

	/*
		NOTE: Using 1.22's routing,
		  GET /posts   -- valid, routes to app.postGetALl
		  GET /posts/  -- INVALID
		  GET /posts/1 -- valid, routes to app.postGet
	*/
	mux.HandleFunc("GET /posts", app.postGetAll)
	mux.HandleFunc("GET /posts/{postID}", app.postGet)
	mux.HandleFunc("POST /posts", app.createPost)
	mux.HandleFunc("PUT /posts/{postID}", app.updatePost)
	mux.HandleFunc("DELETE /posts/{postID}", app.deletePost)
	mux.HandleFunc("GET /posts/{postID}/comments", app.getPostComments)
	mux.HandleFunc("POST /posts/{postID}/comments", app.createPostComment)

	return mux
}
