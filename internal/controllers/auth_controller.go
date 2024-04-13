package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/randnull/banner-service/internal/config"
	"github.com/randnull/banner-service/internal/errors"
	"github.com/randnull/banner-service/internal/jwt_token"
	"github.com/randnull/banner-service/internal/service"
	"github.com/randnull/banner-service/pkg/models"
)


type UserHandlers struct {
	Service *service.UserService
	cfg		*config.Config
}


func NewUserHandler(UserServ *service.UserService, cfg *config.Config) *UserHandlers {
	return &UserHandlers{
		Service:	UserServ,
		cfg:		cfg,
	}
}


func (handler *UserHandlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.Register

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		sendJSONError(w, http.StatusBadRequest, errors.InvalidRequestBody)
		return
	}

	err = handler.Service.CreateUser(user, false)

	if err != nil {
		if err == errors.UsernameAlreadyTaken {
			sendJSONError(w, http.StatusBadRequest, err)
		} else {
			sendJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func (handler *UserHandlers) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user models.Register

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		sendJSONError(w, http.StatusBadRequest, errors.InvalidRequestBody)
		return
	}

	err = handler.Service.CreateUser(user, true)

	if err != nil {
		if err == errors.UsernameAlreadyTaken {
			sendJSONError(w, http.StatusBadRequest, err)
		} else {
			sendJSONError(w, http.StatusInternalServerError, err)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}


func (handler *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user_data models.Register

	err := json.NewDecoder(r.Body).Decode(&user_data)

	if err != nil {
		sendJSONError(w, http.StatusBadRequest, errors.InvalidRequestBody)
		return
	}

	user, err := handler.Service.GetUser(user_data.Username, user_data.Password)

	if err != nil {
		if err == errors.UserNotFound {
			sendJSONError(w, http.StatusNotFound, err)
		} else {
			sendJSONError(w, http.StatusInternalServerError, err)
		}
		return 
	}

	token, _ := jwt_token.CreateJWTToken(user.IsAdmin, handler.cfg.JWTsecret)

	token_json := models.Token{Token: token}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token_json)
}


func (handler *UserHandlers) Auth(hand http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token_str := r.Header.Get("token")

		is_admin, err := jwt_token.ParseJWTToken(token_str, handler.cfg.JWTsecret)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "is_admin", is_admin)

		hand.ServeHTTP(w, r.WithContext(ctx))
	})
}
