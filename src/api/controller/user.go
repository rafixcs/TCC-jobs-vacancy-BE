package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/userfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
)

type CreateUserRequest struct {
	Name     string              `json:"name"`
	Email    string              `json:"email"`
	Password string              `json:"password"`
	RoleId   int                 `json:"role_id"`
	Company  company.CompanyInfo `json:"company"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	userDomain := userfactory.CreateUserDomain()

	err = userDomain.CreateUser(
		userRequest.Name,
		userRequest.Password,
		userRequest.Email,
		userRequest.RoleId,
		userRequest.Company,
	)
	if err != nil {
		message := fmt.Sprintf("failed to create user: %s", err.Error())
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
