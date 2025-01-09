package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/authfactory"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	RoleId int    `json:"role_id"`
}

// Auth godoc
// @Summary authenticate user
// @Description authenticate user
// @Tags Auth
// @Param authrequest body AuthRequest true "Change password"
// @Success 200 {object} AuthResponse "Success"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/auth [post]
func Auth(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	authDomain := authfactory.CreateAuthDomain()

	token, roleId, err := authDomain.UserAuth(authRequest.Email, authRequest.Password)
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

// Logout godoc
// @Summary User logout
// @Description User logout
// @Tags Auth
// @Param Authorization header string true "Authorization token"
// @Success 200 "user logged out"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/user/password [post]
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
