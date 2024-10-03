package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/domain/users"
)

type CreateUserRequest struct {
	Name     string
	Password string
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	err = users.CreateUser(userRequest.Name, userRequest.Password)
	if err != nil {
		message := fmt.Sprintf("failed to create user: %s", err.Error())
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	token, err := users.UserAuth(authRequest.Name, authRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authResponse := AuthResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&authResponse)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	err := users.Logout(tokenHeader)
	if err != nil {
		http.Error(w, "failed to logout", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
