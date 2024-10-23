package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IUserDomain interface {
	CreateUser(name, password, email, companyName, companyEmail, companyLocation, companyDescription string, roleId int) error
}

type UserDomain struct {
	UserRepo    repousers.IUserRepository
	CompanyRepo repocompany.ICompanyRepository
}

func (d *UserDomain) CreateUser(name, password, email, companyName, companyEmail, companyLocation, companyDescription string, roleId int) error {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	if (roleId == 0 || roleId == 1) && companyName == "" {
		return fmt.Errorf("for company user, needs company name")
	}

	alreadyCreatedUser, err := d.UserRepo.CheckIfExists(name)
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

	err = d.UserRepo.Create(userId, name, hashedPassword, email, roleId)

	if roleId == 1 {
		company, err := d.CompanyRepo.FindCompanyByName(companyName)
		if err != nil {
			return err
		}

		if company == (models.Company{}) {
			companyDomain := companyfactory.CreateCompanyDomain()
			companyDomain.CreateCompany(companyName, email, companyDescription, companyLocation)
		}

		err = d.CompanyRepo.CreateUserCompany(company.Id, userId)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}
