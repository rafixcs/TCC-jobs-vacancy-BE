package jobvacancy

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"
)

type IJobVacancyDomain interface {
	CreateJobVacancy(userId, companyId, description, title, location, salary string, requirements, responsabilities []string) error
	CreateUserJobApply(userId, jobId, fullName, email, coverLetter, phone string) error
	GetCompanyJobVacancies(companyName, companyId string) ([]JobVacancyInfo, error)
	GetUserJobApplies(userId string) ([]JobVacancyInfo, error)
	GetUsesAppliesToJobVacancy(jobId string) ([]JobVacancyApplies, error)
	GetJobVacancyDetails(jobId string) (JobVacancyDetails, error)
	SearchJobVacancies(searchStatement string) ([]JobVacancyInfo, error)
}

type JobVacancyDomain struct {
	JobVacancyRepo repojobvacancy.IJobVacancyRepository
	CompanyRepo    repocompany.ICompanyRepository
}

func (d JobVacancyDomain) CreateJobVacancy(
	userId, companyId, description, title, location, salary string,
	requirements, responsabilities []string) error {

	company, err := d.CompanyRepo.FindCompanyById(companyId)
	if err != nil {
		return err
	}

	if company == (models.Company{}) {
		return fmt.Errorf("company not found")
	}

	jsonResp, err := json.Marshal(responsabilities)
	if err != nil {
		return err
	}
	responsabilitiesData := string(jsonResp)

	jsonReqs, err := json.Marshal(requirements)
	if err != nil {
		return err
	}
	requirementsData := string(jsonReqs)

	jobVacancy := models.JobVacancy{
		Id:               uuid.NewString(),
		UserId:           userId,
		CompanyId:        company.Id,
		Title:            title,
		Description:      description,
		Location:         location,
		CreationDate:     time.Now(),
		Responsibilities: responsabilitiesData,
		Requirements:     requirementsData,
		Salary:           salary,
	}

	err = d.JobVacancyRepo.CreateJobVacancy(jobVacancy)
	if err != nil {
		return err
	}

	return nil
}

type JobVacancyDetails struct {
	Id               string
	Description      string    `json:"description"`
	Title            string    `json:"title"`
	Location         string    `json:"location"`
	CreationDate     time.Time `json:"creation_date"`
	Salary           string    `json:"salary"`
	Requirements     []string  `json:"requirements"`
	Responsibilities []string  `json:"responsibilities"`
	Company          string    `json:"company"`
}

func (d JobVacancyDomain) GetJobVacancyDetails(jobId string) (JobVacancyDetails, error) {
	jobVacancyModel, companyName, err := d.JobVacancyRepo.GetJobVacancyDetails(jobId)
	if err != nil {
		return JobVacancyDetails{}, err
	}

	if jobVacancyModel == (models.JobVacancy{}) {
		return JobVacancyDetails{}, fmt.Errorf("job vacancy not found")
	}

	var requirementslist []string
	err = json.Unmarshal([]byte(jobVacancyModel.Requirements), &requirementslist)
	if err != nil {
		return JobVacancyDetails{}, err
	}

	var responsibilitieslist []string
	err = json.Unmarshal([]byte(jobVacancyModel.Responsibilities), &responsibilitieslist)
	if err != nil {
		return JobVacancyDetails{}, err
	}

	jobDetail := JobVacancyDetails{
		Id:               jobVacancyModel.Id,
		Description:      jobVacancyModel.Description,
		Title:            jobVacancyModel.Title,
		CreationDate:     jobVacancyModel.CreationDate,
		Salary:           jobVacancyModel.Salary,
		Requirements:     requirementslist,
		Responsibilities: responsibilitieslist,
		Company:          companyName,
		Location:         jobVacancyModel.Location,
	}

	return jobDetail, nil
}

func (d JobVacancyDomain) CreateUserJobApply(userId, jobId, fullName, email, coverLetter, phone string) error {

	userApply := models.UserApplies{
		Id:           uuid.NewString(),
		UserId:       userId,
		JobVacancyId: jobId,
		FullName:     fullName,
		Email:        email,
		CoverLetter:  coverLetter,
		Phone:        phone,
	}

	err := d.JobVacancyRepo.CreateUserJobApply(userApply)
	if err != nil {
		return err
	}

	return nil
}

func (d JobVacancyDomain) GetCompanyJobVacancies(companyName, companyId string) ([]JobVacancyInfo, error) {
	jobVacanciesModel, err := d.JobVacancyRepo.GetCompanyJobVacancies(companyName, companyId)
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
	UserId      string `json:"user_id"`
	CoverLetter string `json:"cover_letter"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
}

func (d JobVacancyDomain) GetUsesAppliesToJobVacancy(jobId string) ([]JobVacancyApplies, error) {

	usersAppliesModel, err := d.JobVacancyRepo.GetJobVacancyApplies(jobId)
	if err != nil {
		return []JobVacancyApplies{}, err
	}

	var usersApplied []JobVacancyApplies
	for _, model := range usersAppliesModel {
		userApply := JobVacancyApplies{
			UserId:      model.Id,
			CoverLetter: model.CoverLetter,
			Email:       model.Email,
			FullName:    model.FullName,
			Phone:       model.Phone,
		}

		usersApplied = append(usersApplied, userApply)
	}

	return usersApplied, nil
}

func (d JobVacancyDomain) SearchJobVacancies(searchStatement string) ([]JobVacancyInfo, error) {

	jobVacanciesModel, companiesNames, err := d.JobVacancyRepo.SearchJobVacancies(searchStatement)
	if err != nil {
		return []JobVacancyInfo{}, err
	}

	jobVacanciesInfo := JobVacancyInfo{}.TransformSliceModelCompany(jobVacanciesModel, companiesNames)

	return jobVacanciesInfo, nil
}
