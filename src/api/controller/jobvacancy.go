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

	/*var requestContent RegisterUserApplyJobVacancyRequest
	err = json.NewDecoder(r.Body).Decode(&requestContent)
	if err != nil {
		http.Error(w, "bad body format", http.StatusBadRequest)
		return
	}*/

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

func GetCompanyJobVacancies(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	companyId, err := utils.GetCompanyIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
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

func GetUsersAppliesToJobVacancy(w http.ResponseWriter, r *http.Request) {
	jobId := r.URL.Query().Get("job_id")
	if jobId == "" {
		http.Error(w, "missing job_id", http.StatusUnauthorized)
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
