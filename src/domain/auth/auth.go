package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repoauth"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/users"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IAuthDomain interface {
	UserAuth(name, password string) (string, error)
	Logout(tokenHeader string) error
}

type AuthDomain struct {
	AuthRepo repoauth.IAuthRepository
	UserRepo repousers.IUserRepository
}

func (d *AuthDomain) UserAuth(name, password string) (string, error) {
	err := users.UserPasswordValidation(name, password)
	if err != nil {
		return "", err
	}

	userModel, err := d.UserRepo.FindUser(name)
	if err != nil {
		return "", err
	}

	passwordMatching := utils.ValidatePasswordHash(password, userModel.Password)
	if !passwordMatching {
		return "", fmt.Errorf("user/password not matching")
	}

	userLoginModel := models.UserLoginsModel{
		Id:        uuid.NewString(),
		UserId:    userModel.Id,
		LoginDate: time.Now(),
	}

	err = d.AuthRepo.CreateLogin(userLoginModel)
	if err != nil {
		return "", err
	}

	token, err := utils.CreateUserJwtToken(userModel.Id, userLoginModel.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (d *AuthDomain) Logout(tokenHeader string) error {
	token, err := utils.ParseToken(tokenHeader)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	loginId, ok := claims["login_id"].(string)
	if !ok {
		return fmt.Errorf("token missing login field")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return fmt.Errorf("token missing user field")
	}

	validLogin, err := d.AuthRepo.ValidateLogin(loginId, userId)
	if err != nil {
		return err
	}

	if !validLogin {
		return fmt.Errorf("user login not found")
	}

	err = d.AuthRepo.UpdateToLogout(loginId, userId)
	if err != nil {
		return err
	}

	return nil
}
