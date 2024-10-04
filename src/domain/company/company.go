package company

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/repository/models"
)

func CreateCompany(name string) error {
	db, err := datasources.OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()

	companyModel := models.CompanyModels{
		Id:           uuid.NewString(),
		Name:         name,
		CreationDate: time.Now(),
	}

	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM companies WHERE name = $1`
	var alreadyCreated bool
	err = db.QueryRow(query, name).Scan(&alreadyCreated)
	if err != nil {
		return err
	}

	if alreadyCreated {
		return fmt.Errorf("company already created")
	}

	query = `INSERT INTO companies(id, name, creation_date) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, companyModel.Id, companyModel.Name, companyModel.CreationDate)
	if err != nil {
		return err
	}

	return nil
}

func CompaniesList() ([]models.CompanyModels, error) {
	db, err := datasources.OpenDb()
	if err != nil {
		return []models.CompanyModels{}, err
	}
	defer db.Close()

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
