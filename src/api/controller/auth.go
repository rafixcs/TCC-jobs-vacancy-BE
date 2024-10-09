package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repoauth"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/auth"
)

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

	datasource := datasources.DatabasePsql{}
	authRepo := repoauth.AuthRepository{Datasource: &datasource}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	authDomain := auth.AuthDomain{AuthRepo: &authRepo, UserRepo: &userRepo}

	token, err := authDomain.UserAuth(authRequest.Name, authRequest.Password)
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

	datasource := datasources.DatabasePsql{}
	authRepo := repoauth.AuthRepository{Datasource: &datasource}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	authDomain := auth.AuthDomain{AuthRepo: &authRepo, UserRepo: &userRepo}

	err := authDomain.Logout(tokenHeader)
	if err != nil {
		http.Error(w, "failed to logout", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
