package repojobvacancy

import (
	"database/sql"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IJobVacancyRepository interface {
	CreateJobVacancy(jobVacancy models.JobVacancy) error
	CreateUserJobApply(userApply models.UserApplies) error
	GetCompanyJobVacancies(companyName, companyId string) ([]models.JobVacancy, error)
	GetUserJobApplies(userId string) ([]models.JobVacancy, []models.UserApplies, []string, error)
	GetUserJobApply(userId, jobId string) (models.UserApplies, error)
	GetJobVacancyApplies(jobId string) ([]models.UserApplies, error)
	GetJobVacancyDetails(jobId string) (models.JobVacancy, string, error)
	SearchJobVacancies(searchStatement string) ([]models.JobVacancy, []string, error)
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

	query := `INSERT INTO job_vacancies(
			id, user_id, company_id, description,
			title, location, creation_date, responsibilities,
			requirements, salary, job_type, experience_level
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err = db.Exec(
		query,
		jobVacancy.Id,
		jobVacancy.UserId,
		jobVacancy.CompanyId,
		jobVacancy.Description,
		jobVacancy.Title,
		jobVacancy.Location,
		jobVacancy.CreationDate,
		jobVacancy.Responsibilities,
		jobVacancy.Requirements,
		jobVacancy.Salary,
		jobVacancy.JobType,
		jobVacancy.ExperienceLevel,
	)

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

	query := `INSERT INTO user_applies (
				id, job_vacancy_id, user_id, full_name, email, cover_letter, phone, url_resume
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err = db.Exec(query,
		userApply.Id,
		userApply.JobVacancyId,
		userApply.UserId,
		userApply.FullName,
		userApply.Email,
		userApply.CoverLetter,
		userApply.Phone,
		userApply.UrlResume,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *JobVacancyRepository) GetCompanyJobVacancies(companyName, companyId string) ([]models.JobVacancy, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.JobVacancy{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, company_id, user_id, description, title, creation_date, location FROM job_vacancies WHERE company_id IN (SELECT id FROM companies WHERE id = $1 OR name = $2)`
	rows, err := db.Query(query, companyId, companyName)
	if err != nil {
		return []models.JobVacancy{}, err
	}

	var jobVacancies []models.JobVacancy
	for rows.Next() {
		var jobVacancy models.JobVacancy
		rows.Scan(&jobVacancy.Id, &jobVacancy.CompanyId, &jobVacancy.UserId, &jobVacancy.Description, &jobVacancy.Title, &jobVacancy.CreationDate, &jobVacancy.Location)
		jobVacancies = append(jobVacancies, jobVacancy)
	}

	return jobVacancies, nil
}

func (r *JobVacancyRepository) GetJobVacancyApplies(jobId string) ([]models.UserApplies, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.UserApplies{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT id, user_id, job_vacancy_id, full_name, email, phone, cover_letter, url_resume FROM user_applies WHERE job_vacancy_id = $1`
	rows, err := db.Query(query, jobId)
	if err != nil {
		return []models.UserApplies{}, err
	}

	var applies []models.UserApplies
	for rows.Next() {
		var userApply models.UserApplies
		rows.Scan(
			&userApply.Id,
			&userApply.UserId,
			&userApply.JobVacancyId,
			&userApply.FullName,
			&userApply.Email,
			&userApply.Phone,
			&userApply.CoverLetter,
			&userApply.UrlResume,
		)

		applies = append(applies, userApply)
	}

	return applies, nil
}

func (r *JobVacancyRepository) GetUserJobApply(userId, jobId string) (models.UserApplies, error) {

	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.UserApplies{}, nil
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT
				id, job_vacancy_id, user_id, cover_letter, email, full_name, phone, url_resume 
			FROM user_applies
			WHERE user_id = $1 and job_vacancy_id = $2`
	rows, err := db.Query(query, userId, jobId)
	if err != nil {
		return models.UserApplies{}, err
	}

	var userApply models.UserApplies
	if rows.Next() {
		rows.Scan(
			&userApply.Id,
			&userApply.JobVacancyId,
			&userApply.UserId,
			&userApply.CoverLetter,
			&userApply.Email,
			&userApply.FullName,
			&userApply.Phone,
			&userApply.UrlResume,
		)
	}

	return userApply, nil
}

func (r *JobVacancyRepository) GetUserJobApplies(userId string) ([]models.JobVacancy, []models.UserApplies, []string, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.JobVacancy{}, []models.UserApplies{}, []string{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT 
		jv.id, jv.company_id, jv.user_id, jv.description, jv.title, jv.creation_date, cp.name,
		ua.id, ua.cover_letter, ua.email, ua.full_name, ua.phone, ua.url_resume
		FROM job_vacancies AS jv 
		INNER JOIN user_applies AS ua ON ua.job_vacancy_id=jv.id 
		INNER JOIN companies AS cp ON cp.id=jv.company_id
		WHERE ua.user_id = $1`
	rows, err := db.Query(query, userId)
	if err != nil {
		return []models.JobVacancy{}, []models.UserApplies{}, []string{}, err
	}

	var jobVacancies []models.JobVacancy
	var userApplies []models.UserApplies
	var companies []string
	for rows.Next() {
		var jobVacancy models.JobVacancy
		var company string
		var userApply models.UserApplies
		rows.Scan(
			&jobVacancy.Id,
			&jobVacancy.CompanyId,
			&jobVacancy.UserId,
			&jobVacancy.Description,
			&jobVacancy.Title,
			&jobVacancy.CreationDate,
			&company,
			&userApply.Id,
			&userApply.CoverLetter,
			&userApply.Email,
			&userApply.FullName,
			&userApply.Phone,
			&userApply.UrlResume,
		)
		userApply.UserId = jobVacancy.UserId
		userApply.JobVacancyId = jobVacancy.Id

		jobVacancies = append(jobVacancies, jobVacancy)
		companies = append(companies, company)
		userApplies = append(userApplies, userApply)
	}

	return jobVacancies, userApplies, companies, nil
}

func (r *JobVacancyRepository) GetJobVacancyDetails(jobId string) (models.JobVacancy, string, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.JobVacancy{}, "", err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `
		SELECT 
			jv.id, jv.user_id, jv.company_id, jv.description,
			jv.title, jv.creation_date, jv.location, jv.salary,
			jv.requirements, jv.responsibilities, jv.job_type, jv.experience_level,
			cp.name
		FROM job_vacancies AS jv INNER JOIN companies AS cp ON cp.id=jv.company_id WHERE jv.id = $1
		`
	rows, err := db.Query(query, jobId)
	if err != nil {
		return models.JobVacancy{}, "", err
	}

	var jobVacancy models.JobVacancy
	var companyName string
	if rows.Next() {
		rows.Scan(
			&jobVacancy.Id,
			&jobVacancy.UserId,
			&jobVacancy.CompanyId,
			&jobVacancy.Description,
			&jobVacancy.Title,
			&jobVacancy.CreationDate,
			&jobVacancy.Location,
			&jobVacancy.Salary,
			&jobVacancy.Requirements,
			&jobVacancy.Responsibilities,
			&jobVacancy.JobType,
			&jobVacancy.ExperienceLevel,
			&companyName,
		)
	}

	return jobVacancy, companyName, nil
}

func (r *JobVacancyRepository) SearchJobVacancies(searchStatement string) ([]models.JobVacancy, []string, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return []models.JobVacancy{}, []string{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `SELECT 
				jv.id, jv.company_id, jv.user_id, jv.description, jv.title, jv.creation_date, cp.name
			FROM job_vacancies AS jv INNER JOIN companies AS cp ON cp.id=jv.company_id`

	var rows *sql.Rows
	if searchStatement != "" {
		query += ` WHERE search_vector @@ to_tsquery('english', $1)`
		rows, err = db.Query(query, searchStatement)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		return []models.JobVacancy{}, []string{}, err
	}

	var jobVacancies []models.JobVacancy
	var companiesNames []string
	for rows.Next() {
		var jobVacancy models.JobVacancy
		var companyName string
		rows.Scan(
			&jobVacancy.Id,
			&jobVacancy.CompanyId,
			&jobVacancy.UserId,
			&jobVacancy.Description,
			&jobVacancy.Title,
			&jobVacancy.CreationDate,
			&companyName,
		)
		jobVacancies = append(jobVacancies, jobVacancy)
		companiesNames = append(companiesNames, companyName)
	}

	return jobVacancies, companiesNames, nil
}
