package company

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
)

type ICompanyDomain interface {
	CreateCompany(name, email, description, location string) error
	CompaniesList() ([]CompanyInfo, error)
}

type CompanyDomain struct {
	CompanyRepo repocompany.ICompanyRepository
}

func (d *CompanyDomain) CreateCompany(name, email, description, location string) error {
	companyModel := models.Company{
		Id:           uuid.NewString(),
		Name:         name,
		Email:        email,
		Description:  description,
		Location:     location,
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
	if err != nil {
		return err
	}

	return nil
}

type CompanyInfo struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
	Description  string    `json:"description"`
}

func (d *CompanyDomain) CompaniesList() ([]CompanyInfo, error) {
	companiesModels, err := d.CompanyRepo.FindAllCompanies()
	if err != nil {
		return []CompanyInfo{}, err
	}

	var companiesInfo []CompanyInfo
	for _, companyModel := range companiesModels {
		company := CompanyInfo{
			Id:           companyModel.Id,
			Name:         companyModel.Name,
			CreationDate: companyModel.CreationDate,
			Description:  companyModel.Description,
		}

		companiesInfo = append(companiesInfo, company)
	}

	return companiesInfo, nil
}
