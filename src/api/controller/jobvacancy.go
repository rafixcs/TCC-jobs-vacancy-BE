package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/jobfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/jobvacancy"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type CreateJobVacancyRequest struct {
	Description      string   `json:"description"`
	Title            string   `json:"title"`
	Location         string   `json:"location"`
	Salary           string   `json:"salary"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
	JobType          string   `json:"job_type"`
	ExperienceLevel  string   `json:"experience_level"`
}

// CreateJobVacancy godoc
// @Summary Create job vacancy
// @Description Create job vacancy
// @Tags Jobs
// @Param Authorization header string true "Authorization token"
// @Param createjobvacancyrequest body CreateJobVacancyRequest true "Create job vacancy"
// @Success 201 "Created job vacancy"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job [post]
func CreateJobVacancy(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	companyId, err := utils.GetCompanyIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	var requestContent CreateJobVacancyRequest
	json.NewDecoder(r.Body).Decode(&requestContent)

	jobvacancyDomain := jobfactory.CreateJobVacancyDomain()

	err = jobvacancyDomain.CreateJobVacancy(
		userId,
		companyId,
		requestContent.Description,
		requestContent.Title,
		requestContent.Location,
		requestContent.Salary,
		requestContent.JobType,
		requestContent.ExperienceLevel,
		requestContent.Requirements,
		requestContent.Responsibilities,
	)
	if err != nil {
		http.Error(w, "failed to create job vacancy", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetJobVacancyDetails godoc
// @Summary Get job vacancy details
// @Description Get job vacancy details
// @Tags Jobs
// @Param Authorization header string true "Authorization token"
// @Param id path string true "Job vacancy id"
// @Success 200 {object} jobvacancy.JobVacancyDetails "Created job vacancy"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job [post]
func GetJobVacancyDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	jobdomain := jobfactory.CreateJobVacancyDomain()
	jobVacancy, err := jobdomain.GetJobVacancyDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&jobVacancy)
}

type RegisterUserApplyJobVacancyRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	CoverLetter string `json:"cover_letter"`
	JobId       string `json:"job_id"`
	Phone       string `json:"phone"`
}

// RegisterUserApplyJobVacancy godoc
// @Summary Register user apply job vacancy
// @Description Register user apply job vacancy
// @Tags Jobs
// @Param Authorization header string true "Authorization token"
// @Param full_name formData string true "Full name"
// @Param email formData string true "Email"
// @Param phone formData string true "Phone"
// @Param cover_letter formData string true "Cover letter"
// @Param job_id formData string true "Job id"
// @Param resume formData file true "Resume"
// @Success 201 "Created user apply"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job/apply [post]
func RegisterUserApplyJobVacancy(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	var requestContent RegisterUserApplyJobVacancyRequest
	requestContent.FullName = r.FormValue("full_name")
	requestContent.Email = r.FormValue("email")
	requestContent.Phone = r.FormValue("phone")
	requestContent.CoverLetter = r.FormValue("cover_letter")
	requestContent.JobId = r.FormValue("job_id")

	file, handler, err := r.FormFile("resume")
	if err != nil {
		log.Println("Error retrieving file")
		log.Println(err)
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}

	resumeFile := jobvacancy.JobVacancyResumeFile{
		File:   file,
		Header: handler,
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	err = jobVacancyDomain.CreateUserJobApply(
		userId,
		requestContent.JobId,
		requestContent.FullName,
		requestContent.Email,
		requestContent.CoverLetter,
		requestContent.Phone,
		resumeFile,
	)

	if err != nil {
		http.Error(w, "could not create user apply", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type GetCompaniesJobVacanciesResponse struct {
	CompanyName  string                      `json:"company"`
	JobVacancies []jobvacancy.JobVacancyInfo `json:"job_vacancies"`
}

// GetCompanyJobVacancies godoc
// @Summary Get company job vacancies
// @Description Get company job vacancies
// @Tags Jobs
// @Param Authorization header string true "Authorization token"
// @Param company query string false "Company name"
// @Success 200 {object} GetCompaniesJobVacanciesResponse "Company job vacancies"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/company/jobs [get]
func GetCompanyJobVacancies(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	var companyId string
	var err error
	if tokenHeader != "" {
		companyId, err = utils.GetCompanyIdFromToken(tokenHeader)
		if err != nil {
			http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
			return
		}
	}

	companyName := r.URL.Query().Get("company")
	if companyName == "" && companyId == "" {
		http.Error(w, "missing company name/id", http.StatusBadRequest)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobVacancies, err := jobVacancyDomain.GetCompanyJobVacancies(companyName, companyId)
	if err != nil {
		http.Error(w, "failed to get company job vacancies list", http.StatusBadRequest)
		return
	}

	responseBody := GetCompaniesJobVacanciesResponse{
		CompanyName:  companyName,
		JobVacancies: jobVacancies,
	}

	json.NewEncoder(w).Encode(&responseBody)
}

type UserJobAppliesResponse struct {
	JobApplies []jobvacancy.UserJobApply
}

// GetUserJobVacancies godoc
// @Summary Get user job vacancies
// @Description Get user job vacancies
// @Tags Jobs
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} UserJobAppliesResponse "User job vacancies"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job/user [get]
func GetUserJobVacancies(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	userJobApplies, err := jobVacancyDomain.GetUserJobApplies(userId)
	if err != nil {
		http.Error(w, "failed to get user job applies list", http.StatusBadRequest)
		return
	}

	responseBody := UserJobAppliesResponse{
		JobApplies: userJobApplies,
	}

	json.NewEncoder(w).Encode(&responseBody)
}

type SearchJobVacanciesResponse struct {
	JobVacancies []jobvacancy.JobVacancyInfo
}

// SearchJobVacancies godoc
// @Summary Search job vacancies
// @Description Search job vacancies
// @Tags Jobs
// @Param value query string true "Search value"
// @Success 200 {object} SearchJobVacanciesResponse "Job vacancies"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job/search [get]
func SearchJobVacancies(w http.ResponseWriter, r *http.Request) {
	searchStatement := r.URL.Query().Get("value")
	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobVacancies, err := jobVacancyDomain.SearchJobVacancies(searchStatement)
	if err != nil {
		http.Error(w, "failed to search job vacancies", http.StatusBadRequest)
		return
	}

	bodyResponse := SearchJobVacanciesResponse{
		JobVacancies: jobVacancies,
	}

	json.NewEncoder(w).Encode(&bodyResponse)
}

type GetUsersAppliesToJobVacancyResponse struct {
	UsersApplies []jobvacancy.JobVacancyApplies `json:"user_applies"`
}

// GetUsersAppliesToJobVacancy godoc
// @Summary Get users applies to job vacancy
// @Description Get users applies to job vacancy
// @Tags Jobs
// @Param job_id query string true "Job id"
// @Success 200 {object} GetUsersAppliesToJobVacancyResponse "Users applies"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/job/applies [get]
func GetUsersAppliesToJobVacancy(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("job_id")
	if jobId == "" {
		http.Error(w, "missing job id", http.StatusUnauthorized)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	users, err := jobVacancyDomain.GetUsesAppliesToJobVacancy(jobId)
	if err != nil {
		http.Error(w, "failed to get job applies", http.StatusBadRequest)
		return
	}

	responseBody := GetUsersAppliesToJobVacancyResponse{
		UsersApplies: users,
	}

	json.NewEncoder(w).Encode(&responseBody)
}
