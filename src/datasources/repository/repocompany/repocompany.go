package repocompany

import (
	"log"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type ICompanyRepository interface {
	CreateCompany(company models.Company) error
	CreateUserCompany(companyId, userId string) error
	FindCompanyByName(companyName string) (models.Company, error)
	FindCompanyById(companyId string) (models.Company, error)
	FindCompanyByUserId(userId string) (models.Company, error)
	FindIfCompanyExists(companyName string) (bool, error)
	FindAllCompanies() ([]models.Company, error)
}

type CompanyRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *CompanyRepository) CreateCompany(company models.Company) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO companies(id, name, email, description, location, creation_date) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, company.Id, company.Name, company.Email, company.Description, company.Location, company.CreationDate)
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

func (r *CompanyRepository) FindCompanyByName(companyName string) (models.Company, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.Company{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, name, creation_date FROM companies WHERE name = $1`
	row, err := db.Query(query, companyName)
	if err != nil {
		return models.Company{}, err
	}

	if row.Next() {
		var companyModel models.Company
		err = row.Scan(
			&companyModel.Id,
			&companyModel.Name,
			&companyModel.Email,
			&companyModel.CreationDate,
			&companyModel.Description,
			&companyModel.Location,
		)

		if err != nil {
			log.Println(err)
		}

		return companyModel, nil
	}

	return models.Company{}, nil
}

func (r *CompanyRepository) FindCompanyById(companyId string) (models.Company, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.Company{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, name, email, creation_date, description, location FROM companies WHERE id = $1`
	row, err := db.Query(query, companyId)
	if err != nil {
		return models.Company{}, err
	}

	if row.Next() {
		var companyModel models.Company
		err = row.Scan(
			&companyModel.Id,
			&companyModel.Name,
			&companyModel.Email,
			&companyModel.CreationDate,
			&companyModel.Description,
			&companyModel.Location,
		)

		if err != nil {
			log.Println(err)
		}

		return companyModel, nil
	}

	return models.Company{}, nil
}

func (r *CompanyRepository) FindCompanyByUserId(userId string) (models.Company, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.Company{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT 
				cp.id, cp.name, cp.email, cp.creation_date, cp.description, cp.location 
			  	FROM companies AS cp 
				INNER JOIN company_users AS cpu ON cpu.company_id=cp.id
				WHERE cpu.user_id = $1`

	row, err := db.Query(query, userId)
	if err != nil {
		return models.Company{}, err
	}

	if row.Next() {
		var companyModel models.Company
		err = row.Scan(
			&companyModel.Id,
			&companyModel.Name,
			&companyModel.Email,
			&companyModel.CreationDate,
			&companyModel.Description,
			&companyModel.Location,
		)

		if err != nil {
			return models.Company{}, err
		}

		return companyModel, nil
	}

	return models.Company{}, nil
}

func (r *CompanyRepository) FindAllCompanies() ([]models.Company, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.Company{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var companies []models.Company

	query := `SELECT id, name, description, creation_date FROM companies`
	rows, err := db.Query(query)
	if err != nil {
		return []models.Company{}, err
	}

	for rows.Next() {
		var companyModel models.Company
		err = rows.Scan(&companyModel.Id, &companyModel.Name, &companyModel.Description, &companyModel.CreationDate)
		if err != nil {
			return []models.Company{}, nil
		}
		companies = append(companies, companyModel)
	}

	return companies, nil
}
