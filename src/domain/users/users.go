package users

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IUserDomain interface {
	CreateUser(name, password, email, phone string, roleId int, company company.CompanyInfo) error
	UserDetails(userId string) (UserDetails, error)
	UpdateUser(userId, name, phone string) error
	ChangePassword(userId, oldPassword, newPassword string) error
}

type UserDomain struct {
	UserRepo    repousers.IUserRepository
	CompanyRepo repocompany.ICompanyRepository
}

func (d *UserDomain) CreateUser(name, password, email, phone string, roleId int, company company.CompanyInfo) error {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	if (roleId == 1) && company.Name == "" {
		return fmt.Errorf("for company user, needs company name")
	}

	alreadyCreatedUser, err := d.UserRepo.CheckIfExists(name, email)
	if err != nil {
		return err
	}

	if alreadyCreatedUser {
		return fmt.Errorf(`user already created`)
	}

	userId := uuid.NewString()

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = d.UserRepo.Create(userId, name, hashedPassword, email, phone, roleId)
	if err != nil {
		return err
	}

	err = d.handleCompanyUserCreate(userId, roleId, company)
	if err != nil {
		return err
	}

	return nil
}

func (d *UserDomain) handleCompanyUserCreate(userId string, roleId int, companyInfo company.CompanyInfo) error {
	if roleId == 1 {
		company, err := d.CompanyRepo.FindCompanyByName(companyInfo.Name)
		if err != nil {
			return err
		}

		if company == (models.Company{}) {
			companyDomain := companyfactory.CreateCompanyDomain()
			company, err = companyDomain.CreateCompany(companyInfo.Name, companyInfo.Email, companyInfo.Description, companyInfo.Location)
			if err != nil {
				if err.Error() != "company already created" {
					return err
				}
			}
		}

		err = d.CompanyRepo.CreateUserCompany(company.Id, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

type UserDetails struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	RoleId int    `json:"role_id"`
}

func (d *UserDomain) UserDetails(userId string) (UserDetails, error) {
	userModel, err := d.UserRepo.FindUserById(userId)
	if err != nil {
		return UserDetails{}, err
	}

	userDetails := UserDetails{
		Name:   userModel.Name,
		Email:  userModel.Email,
		Phone:  userModel.Phone,
		RoleId: userModel.RoleId,
	}

	return userDetails, nil
}

func (d *UserDomain) UpdateUser(userId, name, phone string) error {

	userModel := models.User{
		Id:    userId,
		Name:  name,
		Phone: phone,
	}

	err := d.UserRepo.UpdateUser(userModel)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("couldn't update user")
	}

	return nil
}

func (d *UserDomain) ChangePassword(userId, oldPassword, newPassword string) error {
	userModel, err := d.UserRepo.FindUserById(userId)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("couldn't find user")
	}

	passwordMatching := utils.ValidatePasswordHash(oldPassword, userModel.Password)
	if !passwordMatching {
		return fmt.Errorf("password not matching")
	}

	newHashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to hash password")
	}

	userModel.Password = newHashedPassword

	err = d.UserRepo.UpdatePassword(userModel)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("couldn't update password")
	}

	return nil
}
