package repojobvacancy

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IJobVacancyRepository interface {
	CreateJobVacancy(jobVacancy models.JobVacancy) error
	CreateUserJobApply(userApply models.UserApplies) error
	GetCompanyJobVacancies(companyName string) ([]models.JobVacancy, error)
	GetUserJobApplies(userId string) ([]models.JobVacancy, error)
}

type JobVacancyRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *JobVacancyRepository) CreateJobVacancy(jobVacancy models.JobVacancy) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := "INSERT INTO job_vacancies(id, user_id, company_id, description, title, creation_date) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = db.Exec(query, jobVacancy.Id, jobVacancy.UserId, jobVacancy.CompanyId, jobVacancy.Description, jobVacancy.Title, jobVacancy.CreationDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *JobVacancyRepository) CreateUserJobApply(userApply models.UserApplies) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO user_applies (id, job_vacancy_id, user_id) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, userApply.Id, userApply.JobVacancyId, userApply.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (r *JobVacancyRepository) GetCompanyJobVacancies(companyName string) ([]models.JobVacancy, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.JobVacancy{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, company_id, user_id, description, title, creation_date FROM job_vacancies WHERE company_id IN (SELECT id FROM companies WHERE name = $1)`
	rows, err := db.Query(query, companyName)
	if err != nil {
		return []models.JobVacancy{}, err
	}

	var jobVacancies []models.JobVacancy
	for rows.Next() {
		var jobVacancy models.JobVacancy
		rows.Scan(&jobVacancy.Id, &jobVacancy.CompanyId, &jobVacancy.UserId, &jobVacancy.Description, &jobVacancy.Title, &jobVacancy.CreationDate)
		jobVacancies = append(jobVacancies, jobVacancy)
	}

	return jobVacancies, nil

}

func (r *JobVacancyRepository) GetUserJobApplies(userId string) ([]models.JobVacancy, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.JobVacancy{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, company_id, user_id, description, title, creation_date FROM job_vacancies WHERE user_id = $1`
	rows, err := db.Query(query, userId)
	if err != nil {
		return []models.JobVacancy{}, err
	}

	var jobVacancies []models.JobVacancy
	for rows.Next() {
		var jobVacancy models.JobVacancy
		rows.Scan(&jobVacancy.Id, &jobVacancy.CompanyId, &jobVacancy.UserId, &jobVacancy.Description, &jobVacancy.Title, &jobVacancy.CreationDate)
		jobVacancies = append(jobVacancies, jobVacancy)
	}

	return jobVacancies, nil
}
