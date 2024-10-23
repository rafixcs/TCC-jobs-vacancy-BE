package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/userfactory"
)

type CreateUserRequest struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	CompanyName        string `json:"company"`
	CompanyEmail       string `json:"company_email"`
	CompanyLocation    string `json:"company_location"`
	CompanyDescription string `json:"company_description"`
	RoleId             int    `json:"role_id"`
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
		userRequest.CompanyName,
		userRequest.CompanyEmail,
		userRequest.CompanyLocation,
		userRequest.CompanyDescription,
		userRequest.RoleId,
	)
	if err != nil {
		message := fmt.Sprintf("failed to create user: %s", err.Error())
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
