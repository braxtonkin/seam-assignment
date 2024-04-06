package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/braxtonkin/blogapi/internal/data"
	"net/http"
	"strconv"
)

// Handler for retrieving a single blog by its id
func (app *application) postGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process postID")
		return
	}

	blog, err := app.blogModel.Get(id)
	if err != nil {
		if errors.Is(err, data.RecordNotFound) {
			app.writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			app.writeServerError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, blog)
	if err != nil {
		app.writeServerError(w, r, err)
	}
}

// Handler returning all blog posts
func (app *application) postGetAll(w http.ResponseWriter, r *http.Request) {
	blogs, err := app.blogModel.GetAll()
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	if len(blogs) == 0 {
		app.writeJSON(w, http.StatusOK, wrapper{"": "no records found"})
		return
	}
	app.writeJSON(w, http.StatusOK, blogs)

}

// Handler for creating a blog post
func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string
		Content string
		Author  string
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process JSON body")
		return
	}

	blog := &data.Blog{
		Title:   input.Title,
		Content: input.Content,
		Author:  input.Author,
	}

	id, err := app.blogModel.Insert(blog)
	if err != nil {
		app.writeServerError(w, r, err)
	}

	// Set a header to tell the client where the new endpoint is located
	w.Header().Set("Location", fmt.Sprintf("/posts/%d", id))

	err = app.writeJSON(w, http.StatusCreated, wrapper{"id": id, "blog": blog})
	if err != nil {
		app.writeServerError(w, r, err)
	}
}

// Handler for updating an existing blog post
func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process postID")
		return
	}

	blog, err := app.blogModel.Get(id)
	if err != nil {
		if errors.Is(err, data.RecordNotFound) {
			app.writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			app.writeServerError(w, r, err)
		}
		return
	}

	var input struct {
		Title   string
		Content string
		Author  string
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process JSON body")
		return
	}

	blog.Title = input.Title
	blog.Author = input.Author
	blog.Content = input.Content

	err = app.blogModel.Update(&blog, id)
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, blog)
	if err != nil {
		app.writeServerError(w, r, err)
	}
}

// Handler for deleting an existing blog post
func (app *application) deletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process postID")
		return
	}

	err = app.blogModel.Delete(id)
	if err != nil {
		if errors.Is(err, data.RecordNotFound) {
			app.writeError(w, r, http.StatusNotFound, "unable to find blog")
		} else {
			app.writeServerError(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, wrapper{"message": "successfully deleted"})
	if err != nil {
		app.writeServerError(w, r, err)
	}
}

// Handler for getting all comments for the blog post associated with the postID
func (app *application) getPostComments(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process postID")
		return
	}

	comments, err := app.commentModel.Get(id)
	if err != nil {
		app.writeServerError(w, r, err)
		return
	}

	if len(comments) == 0 {
		app.writeJSON(w, http.StatusOK, wrapper{"": "no records found"})
		return
	}
	app.writeJSON(w, http.StatusOK, comments)
}

// Handler for creating a new comment
func (app *application) createPostComment(w http.ResponseWriter, r *http.Request) {
	blogId, err := strconv.Atoi(r.PathValue("postID"))
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process postID")
		return
	}

	var input struct {
		Content string
		Author  string
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.writeError(w, r, http.StatusBadRequest, "unable to process JSON body")
		return
	}

	comment := &data.Comment{
		Content: input.Content,
		Author:  input.Author,
	}

	commentId, err := app.commentModel.Insert(comment, blogId)
	if err != nil {
		if errors.Is(err, data.ForeignKeyError) {
			app.writeError(w, r, http.StatusNotFound, err.Error())
		} else {
			app.writeServerError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, wrapper{"id": commentId, "comment": comment})
	if err != nil {
		app.writeServerError(w, r, err)
	}
}
