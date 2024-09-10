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
	}

	w.WriteHeader(http.StatusCreated)
}
