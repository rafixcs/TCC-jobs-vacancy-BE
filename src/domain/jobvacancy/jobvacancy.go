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
