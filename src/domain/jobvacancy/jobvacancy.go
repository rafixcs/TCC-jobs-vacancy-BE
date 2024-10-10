package jobvacancy

import (
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"
)

type IJobVacancyDomain interface {
	CreateJobVacancy(userId, companyName, description, title string) error
	CreateUserJobApply(userId, jobId string) error
	GetCompanyJobVacancies(companyName string) ([]JobVacancy, error)
	GetUserJobApplies(userId string) ([]JobVacancy, error)
}

type JobVacancyDomain struct {
	JobVacancyRepo repojobvacancy.IJobVacancyRepository
	CompanyRepo    repocompany.ICompanyRepository
}

func (d JobVacancyDomain) CreateJobVacancy(userId, companyName, description, title string) error {
	company, err := d.CompanyRepo.FindCompanyByName(companyName)
	if err != nil {
		return err
	}

	jobVacancy := models.JobVacancy{
		Id:           uuid.NewString(),
		UserId:       userId,
		CompanyId:    company.Id,
		Title:        title,
		Description:  description,
		CreationDate: time.Now(),
	}

	err = d.JobVacancyRepo.CreateJobVacancy(jobVacancy)
	if err != nil {
		return err
	}

	return nil
}

func (d JobVacancyDomain) CreateUserJobApply(userId, jobId string) error {

	userApply := models.UserApplies{
		Id:           uuid.NewString(),
		UserId:       userId,
		JobVacancyId: jobId,
	}

	err := d.JobVacancyRepo.CreateUserJobApply(userApply)
	if err != nil {
		return err
	}

	return nil
}

type JobVacancy struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
}

func (d JobVacancyDomain) GetCompanyJobVacancies(companyName string) ([]JobVacancy, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.GetCompanyJobVacancies(companyName)
	if err != nil {
		return []JobVacancy{}, err
	}

	var jobVacanciesInfo []JobVacancy

	for _, job := range jobVacanciesModel {
		info := JobVacancy{
			Title:        job.Title,
			Description:  job.Description,
			CreationDate: job.CreationDate,
		}
		jobVacanciesInfo = append(jobVacanciesInfo, info)
	}

	return jobVacanciesInfo, nil
}

func (d JobVacancyDomain) GetUserJobApplies(userId string) ([]JobVacancy, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.GetUserJobApplies(userId)
	if err != nil {
		return []JobVacancy{}, err
	}

	var jobVacanciesInfo []JobVacancy

	for _, job := range jobVacanciesModel {
		info := JobVacancy{
			Title:        job.Title,
			Description:  job.Description,
			CreationDate: job.CreationDate,
		}
		jobVacanciesInfo = append(jobVacanciesInfo, info)
	}

	return jobVacanciesInfo, nil
}
