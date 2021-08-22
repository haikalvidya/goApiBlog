package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/haikalvidya/goApiBlog/api/auth"
	"github.com/haikalvidya/goApiBlog/api/models"
	"github.com/haikalvidya/goApiBlog/api/responses"
	"github.com/haikalvidya/goApiBlog/api/formaterror"
)

func (server *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	// parsing body request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// get parameter from request body
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// initialize post
	post.Init()
	err = post.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// is user authenticate?
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != post.AuthorId {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formatedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, postCreated.ID))
	responses.JSON(w, http.StatusCreated, postCreated)
}

func (server *Server) GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}

func (server *Server) GetPost(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	pid, err := strconv.ParseUint(theVars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}
	postReceived, err := post.FindPostById(server.DB, pid)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, postReceived)
}

func (server *Server) UpdatePost(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	// check if the post id is valid
	pid, err := strconv.ParseUint(theVars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// check if auth token is valid and get the user_id from it
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.Error(w, http.StatusNotFound, errors.New("Post Not Foudn"))
		return
	}

	// if a user attempt to update post with no belonging to him
	if uid != post.AuthorId {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// read data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// start processing request data
	postUpdate := models.Post{}
	err = json.Unmarshal(body, &postUpdate)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate.Init()
	err = postUpdate.Validate()
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	// update post_id
	postUpdate.ID = post.ID

	postUpdated, err := postUpdate.UpdateAPost(server.DB)

	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formatedError)
		return
	}
	responses.JSON(w, http.StatusOK, postUpdated)
}

func (server *Server) DeletePost(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	// check, is valid post_id that given
	pid, err := strconv.ParseUint(theVars["id"], 10, 64)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}

	// check user authorization
	uid, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// check existance of post
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.Error(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// check is user the owner of the post/author?
	if uid != post.AuthorId {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Deleting post
	_, err = post.DeleteAPost(server.DB, pid, uid)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}