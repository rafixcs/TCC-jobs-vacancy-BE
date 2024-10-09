package repojobvacancy

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IJobVacancyRepository interface {
	CreateJobVacancy(jobVacancy models.JobVacancy) error
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
