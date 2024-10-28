package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repoauth"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/users"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IAuthDomain interface {
	UserAuth(email, password string) (string, int, error)
	Logout(tokenHeader string) error
}

type AuthDomain struct {
	AuthRepo repoauth.IAuthRepository
	UserRepo repousers.IUserRepository
}

func (d *AuthDomain) UserAuth(email, password string) (string, int, error) {
	err := users.UserPasswordValidation(email, password)
	if err != nil {
		return "", -1, err
	}

	userModel, err := d.UserRepo.FindUserByEmail(email)
	if err != nil {
		return "", -1, err
	}

	passwordMatching := utils.ValidatePasswordHash(password, userModel.Password)
	if !passwordMatching {
		return "", -1, fmt.Errorf("user/password not matching")
	}

	userLoginModel := models.UserLogins{
		Id:        uuid.NewString(),
		UserId:    userModel.Id,
		LoginDate: time.Now(),
	}

	err = d.AuthRepo.CreateLogin(userLoginModel)
	if err != nil {
		return "", -1, err
	}

	var token string
	if userModel.RoleId == 1 {
		companyDomain := companyfactory.CreateCompanyDomain()
		companyInfo, err := companyDomain.GetUserCompany(userModel.Id)
		if err != nil {
			return "", -1, err
		}

		token, err = utils.CreateUserCompanyJwtToken(userModel.Id, userLoginModel.Id, companyInfo.Id)
	} else {
		token, err = utils.CreateUserJwtToken(userModel.Id, userLoginModel.Id)
	}

	if err != nil {
		return "", -1, err
	}

	return token, userModel.RoleId, nil
}

func (d *AuthDomain) Logout(tokenHeader string) error {
	userId, loginId, err := utils.GetUserAuthIdsFromToken(tokenHeader)
	if err != nil {
		return err
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
