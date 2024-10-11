package jobvacancy

import (
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"
)

type IJobVacancyDomain interface {
	CreateJobVacancy(userId, companyName, description, title, location string) error
	CreateUserJobApply(userId, jobId string) error
	GetCompanyJobVacancies(companyName string) ([]JobVacancyInfo, error)
	GetUserJobApplies(userId string) ([]JobVacancyInfo, error)
	GetUsesAppliesToJobVacancy(jobId string) ([]JobVacancyApplies, error)
	SearchJobVacancies(searchStatement string) ([]JobVacancyInfo, error)
}

type JobVacancyDomain struct {
	JobVacancyRepo repojobvacancy.IJobVacancyRepository
	CompanyRepo    repocompany.ICompanyRepository
}

func (d JobVacancyDomain) CreateJobVacancy(userId, companyName, description, title, location string) error {
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
		Location:     location,
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

func (d JobVacancyDomain) GetCompanyJobVacancies(companyName string) ([]JobVacancyInfo, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.GetCompanyJobVacancies(companyName)
	if err != nil {
		return []JobVacancyInfo{}, err
	}

	jobVacanciesInfo := JobVacancyInfo{}.TransformSliceModel(jobVacanciesModel)

	return jobVacanciesInfo, nil
}

func (d JobVacancyDomain) GetUserJobApplies(userId string) ([]JobVacancyInfo, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.GetUserJobApplies(userId)
	if err != nil {
		return []JobVacancyInfo{}, err
	}

	jobVacanciesInfo := JobVacancyInfo{}.TransformSliceModel(jobVacanciesModel)

	return jobVacanciesInfo, nil
}

type JobVacancyApplies struct {
	UserId   string
	UserName string
}

func (d JobVacancyDomain) GetUsesAppliesToJobVacancy(jobId string) ([]JobVacancyApplies, error) {

	usersModels, err := d.JobVacancyRepo.GetJobVacancyApplies(jobId)
	if err != nil {
		return []JobVacancyApplies{}, err
	}

	var usersApplied []JobVacancyApplies
	for _, model := range usersModels {
		userApply := JobVacancyApplies{
			UserId:   model.Id,
			UserName: model.Name,
		}

		usersApplied = append(usersApplied, userApply)
	}

	return usersApplied, nil
}

func (d JobVacancyDomain) SearchJobVacancies(searchStatement string) ([]JobVacancyInfo, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.SearchJobVacancies(searchStatement)
	if err != nil {
		return []JobVacancyInfo{}, err
	}

	jobVacanciesInfo := JobVacancyInfo{}.TransformSliceModel(jobVacanciesModel)

	return jobVacanciesInfo, nil
}
