package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/authfactory"
)

type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	RoleId int    `json:"role_id"`
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	authDomain := authfactory.CreateAuthDomain()

	token, roleId, err := authDomain.UserAuth(authRequest.Name, authRequest.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	authResponse := AuthResponse{
		Token:  token,
		RoleId: roleId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&authResponse)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")

	authDomain := authfactory.CreateAuthDomain()

	err := authDomain.Logout(tokenHeader)
	if err != nil {
		http.Error(w, "failed to logout", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
