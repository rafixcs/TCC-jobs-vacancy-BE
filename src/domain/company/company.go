package company

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
)

type ICompanyDomain interface {
	CreateCompany(name string) error
	CompaniesList() ([]models.CompanyModels, error)
}

type CompanyDomain struct {
	CompanyRepo repocompany.ICompanyRepository
}

func (d *CompanyDomain) CreateCompany(name string) error {
	companyModel := models.CompanyModels{
		Id:           uuid.NewString(),
		Name:         name,
		CreationDate: time.Now(),
	}

	alreadyCreated, err := d.CompanyRepo.FindIfCompanyExists(name)
	if err != nil {
		return err
	}

	if alreadyCreated {
		return fmt.Errorf("company already created")
	}

	err = d.CompanyRepo.CreateCompany(companyModel)

	return nil
}

func (d *CompanyDomain) CompaniesList() ([]models.CompanyModels, error) {
	return d.CompanyRepo.FindAllCompanies()
}
