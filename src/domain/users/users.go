package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IUserDomain interface {
	CreateUser(name, password, email string, roleId int, company company.CompanyInfo) error
}

type UserDomain struct {
	UserRepo    repousers.IUserRepository
	CompanyRepo repocompany.ICompanyRepository
}

func (d *UserDomain) CreateUser(name, password, email string, roleId int, company company.CompanyInfo) error {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	if (roleId == 1) && company.Name == "" {
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
				return err
			}
		}

		err = d.CompanyRepo.CreateUserCompany(company.Id, userId)
		if err != nil {
			return err
		}
	}

	return nil
}
