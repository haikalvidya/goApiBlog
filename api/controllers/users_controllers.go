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

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Init()
	err = user.Validate("")
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	userCreated, err := user.SaveUser(server.DB)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formatedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	// parsing parameter from request
	theVars := mux.Vars(r)
	uid, err := strconv.ParseUint(theVars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	// use find user by id method from models user
	user := models.User{}
	theUser, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
	}
	responses.JSON(w, http.StatusOK, theUser)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// parsing paramater
	theVars := mux.Vars(r)
	uid, err := strconv.ParseUint(theVars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	// get the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	// get user models
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	// get token id
	tokenID, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	// proccess updating
	user.Init()
	err = user.Validate("update")
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	updateUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.Error(w, http.StatusInternalServerError, formatedError)
		return
	}
	responses.JSON(w, http.StatusOK, updateUser)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	theVars := mux.Vars(r)
	user := models.User{}

	uid, err := strconv.ParseUint(theVars["id"], 10, 32)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenId(r)
	if err != nil {
		responses.Error(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		responses.Error(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")
}