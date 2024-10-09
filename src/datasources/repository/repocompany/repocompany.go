package repocompany

import (
	"log"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/models"
)

type ICompanyRepository interface {
	CreateCompany(company models.CompanyModels) error
	CreateUserCompany(companyId, userId string) error
	FindCompanyByName(companyName string) (models.CompanyModels, error)
	FindIfCompanyExists(companyName string) (bool, error)
	FindAllCompanies() ([]models.CompanyModels, error)
}

type CompanyRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *CompanyRepository) CreateCompany(company models.CompanyModels) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO companies(id, name, creation_date) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, company.Id, company.Name, company.CreationDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) CreateUserCompany(companyId, userId string) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO company_users(company_id, user_id) VALUES ($1, $2)`
	_, err = db.Exec(query, companyId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *CompanyRepository) FindIfCompanyExists(companyName string) (bool, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return false, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM companies WHERE name = $1`
	var alreadyCreated bool
	err = db.QueryRow(query, companyName).Scan(&alreadyCreated)
	if err != nil {
		return false, err
	}

	return alreadyCreated, nil
}

func (r *CompanyRepository) FindCompanyByName(companyName string) (models.CompanyModels, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.CompanyModels{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, name, creation_date FROM companies WHERE name = $1`
	row, err := db.Query(query, companyName)
	if err != nil {
		return models.CompanyModels{}, err
	}

	if row.Next() {
		var companyModel models.CompanyModels
		err = row.Scan(&companyModel.Id, &companyModel.Name, &companyModel.CreationDate)
		if err != nil {
			log.Println(err)
		}

		return companyModel, nil
	}

	return models.CompanyModels{}, nil
}

func (r *CompanyRepository) FindAllCompanies() ([]models.CompanyModels, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.CompanyModels{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var companies []models.CompanyModels

	query := `SELECT id, name, creation_date FROM companies`
	rows, err := db.Query(query)
	if err != nil {
		return []models.CompanyModels{}, err
	}

	for rows.Next() {
		var companyModel models.CompanyModels
		err = rows.Scan(&companyModel.Id, &companyModel.Name, &companyModel.CreationDate)
		if err != nil {
			log.Println(err)
		}
		companies = append(companies, companyModel)
	}

	return companies, nil
}
