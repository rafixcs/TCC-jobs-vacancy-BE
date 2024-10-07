package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/users"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	RoleId   int    `json:"role_id"`
}

// Add factory
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	datasource := datasources.DatabasePsql{}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	userDomain := users.UserDomain{UserRepo: &userRepo}

	err = userDomain.CreateUser(userRequest.Name, userRequest.Password, userRequest.RoleId)
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
