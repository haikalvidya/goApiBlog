package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/haikalvidya/goApiBlog/api/auth"
	"github.com/haikalvidya/goApiBlog/api/models"
	"github.com/haikalvidya/goApiBlog/api/responses"
	"github.com/haikalvidya/goApiBlog/api/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Init()
	err = user.Validate("login")
	if err != nil {
		responses.Error(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formatedError := formaterror.FormatError(err.Error())
		responses.Error(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {
	user := models.User{}
	err := server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassowrd(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}