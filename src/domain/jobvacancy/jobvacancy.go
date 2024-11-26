package jobvacancy

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	config "github.com/rafixcs/tcc-job-vacancy/src/configuration"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type IJobVacancyDomain interface {
	CreateJobVacancy(userId, companyId, description, title, location, salary, jobType, experienceLevel string, requirements, responsabilities []string) error
	CreateUserJobApply(userId, jobId, fullName, email, coverLetter, phone string, resumeFile JobVacancyResumeFile) error
	GetCompanyJobVacancies(companyName, companyId string) ([]JobVacancyInfo, error)
	GetUserJobApplies(userId string) ([]UserJobApply, error)
	GetUsesAppliesToJobVacancy(jobId string) ([]JobVacancyApplies, error)
	GetJobVacancyDetails(jobId string) (JobVacancyDetails, error)
	SearchJobVacancies(searchStatement string) ([]JobVacancyInfo, error)
}

type JobVacancyDomain struct {
	JobVacancyRepo repojobvacancy.IJobVacancyRepository
	CompanyRepo    repocompany.ICompanyRepository
}

func (d JobVacancyDomain) CreateJobVacancy(
	userId, companyId, description, title, location, salary, jobType, experienceLevel string,
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
		JobType:          jobType,
		ExperienceLevel:  experienceLevel,
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
	JobType          string    `json:"job_type"`
	ExperienceLevel  string    `json:"experience_level"`
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
		JobType:          jobVacancyModel.JobType,
		ExperienceLevel:  jobVacancyModel.ExperienceLevel,
	}

	return jobDetail, nil
}

func (d JobVacancyDomain) CreateUserJobApply(userId, jobId, fullName, email, coverLetter, phone string, resumeFile JobVacancyResumeFile) error {

	defer resumeFile.File.Close()

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.CF_ACCESS_KEY, config.CF_SECRET_ACCESS_KEY, ""),
		Endpoint:    aws.String(config.R2_ENDPOINT),
		Region:      aws.String("us-east-1"),
	})

	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error creating AWS session")
	}

	svc := s3.New(sess)
	objectKey := fmt.Sprintf("%d-%s", time.Now().UnixNano(), resumeFile.Header.Filename)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.R2_BUCKET),
		Key:    aws.String(objectKey),
		Body:   resumeFile.File,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error uploading file")
	}

	fileURL := fmt.Sprintf("%s/%s/%s", config.R2_ENDPOINT, config.R2_BUCKET, objectKey)

	tempFile, err := os.Create("/app/uploads/" + resumeFile.Header.Filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resumeFile.File)
	if err != nil {
		log.Println(err)
		return err
	}

	userApply := models.UserApplies{
		Id:           uuid.NewString(),
		UserId:       userId,
		JobVacancyId: jobId,
		FullName:     fullName,
		Email:        email,
		CoverLetter:  coverLetter,
		Phone:        phone,
		UrlResume:    fileURL,
	}

	err = d.JobVacancyRepo.CreateUserJobApply(userApply)
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

type UserJobApply struct {
	JobInfo   JobVacancyInfo
	UserApply models.UserApplies
}

func (d JobVacancyDomain) GetUserJobApplies(userId string) ([]UserJobApply, error) {
	jobVacanciesModel, userAppliesModel, companies, err := d.JobVacancyRepo.GetUserJobApplies(userId)
	if err != nil {
		return []UserJobApply{}, err
	}

	jobVacanciesInfo := JobVacancyInfo{}.TransformSliceModelCompany(jobVacanciesModel, companies)

	var userAppliesJob []UserJobApply
	for i, job := range jobVacanciesInfo {
		userApplyJob := UserJobApply{
			JobInfo:   job,
			UserApply: userAppliesModel[i],
		}
		userAppliesJob = append(userAppliesJob, userApplyJob)
	}

	return userAppliesJob, nil
}

type JobVacancyApplies struct {
	UserId      string `json:"user_id"`
	CoverLetter string `json:"cover_letter"`
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone"`
	UrlResume   string `json:"url_resume"`
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
			UrlResume:   model.UrlResume,
		}

		sess, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(config.CF_ACCESS_KEY, config.CF_SECRET_ACCESS_KEY, ""),
			Endpoint:    aws.String(config.R2_ENDPOINT),
			Region:      aws.String("us-east-1"),
		})

		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("error creating AWS session")
		}
		svc := s3.New(sess)
		segments := strings.Split(userApply.UrlResume, "/")
		objectKey := segments[len(segments)-1]
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(config.R2_BUCKET),
			Key:    aws.String(objectKey),
		})

		urlStr, err := req.Presign(15 * time.Minute)
		if err != nil {
			return nil, err
		}
		userApply.UrlResume = urlStr

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
