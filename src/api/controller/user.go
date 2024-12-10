package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/userfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"

	_ "github.com/rafixcs/tcc-job-vacancy/src/docs"
)

type CreateUserRequest struct {
	Name     string              `json:"name"`
	Email    string              `json:"email"`
	Password string              `json:"password"`
	RoleId   int                 `json:"role_id"`
	Company  company.CompanyInfo `json:"company"`
	Phone    string              `json:"phone"`
}

// CreateUser godoc
// @Summary Create user
// @Description Create user
// @Tags User
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/v1/user [post]
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
		userRequest.Phone,
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

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "bad token request", http.StatusBadRequest)
		return
	}

	userDomain := userfactory.CreateUserDomain()
	userDetails, err := userDomain.UserDetails(userId)
	if err != nil {
		http.Error(w, "couldnt get user details", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(userDetails)
	w.WriteHeader(http.StatusOK)
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "bad token request", http.StatusBadRequest)
		return
	}

	var updateUserRequest UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	userDomain := userfactory.CreateUserDomain()
	err = userDomain.UpdateUser(userId, updateUserRequest.Name, updateUserRequest.Phone)
	if err != nil {
		http.Error(w, "couldn't update user", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "bad token request", http.StatusBadRequest)
		return
	}

	var changePasswordRequest ChangePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&changePasswordRequest)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}

	userDomain := userfactory.CreateUserDomain()
	err = userDomain.ChangePassword(userId, changePasswordRequest.OldPassword, changePasswordRequest.NewPassword)
	if err != nil {
		http.Error(w, "couldn't update password", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
