package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/jobfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/jobvacancy"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type CreateJobVacancyRequest struct {
	CompanyName      string   `json:"company"`
	Description      string   `json:"description"`
	Title            string   `json:"title"`
	Location         string   `json:"location"`
	Salary           string   `json:"salary"`
	Requirements     []string `json:"requirements"`
	Responsibilities []string `json:"responsibilities"`
}

func CreateJobVacancy(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	var requestContent CreateJobVacancyRequest
	json.NewDecoder(r.Body).Decode(&requestContent)

	jobvacancyDomain := jobfactory.CreateJobVacancyDomain()

	err = jobvacancyDomain.CreateJobVacancy(
		userId,
		requestContent.CompanyName,
		requestContent.Description,
		requestContent.Title,
		requestContent.Location,
		requestContent.Salary,
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

func RegisterUserApplyJobVacancy(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	jobId := r.URL.Query().Get("job_id")
	if jobId == "" {
		http.Error(w, "missing job_id", http.StatusUnauthorized)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()

	err = jobVacancyDomain.CreateUserJobApply(userId, jobId)
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
	companyName := r.URL.Query().Get("company")
	if companyName == "" {
		http.Error(w, "missing company name", http.StatusBadRequest)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobVacancies, err := jobVacancyDomain.GetCompanyJobVacancies(companyName)
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
	JobApplies []jobvacancy.JobVacancyInfo
}

func GetUserJobVacancies(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	jobVacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobApplies, err := jobVacancyDomain.GetUserJobApplies(userId)
	if err != nil {
		http.Error(w, "failed to get user job applies list", http.StatusBadRequest)
		return
	}

	responseBody := UserJobAppliesResponse{
		JobApplies: jobApplies,
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
	UsersApplies []jobvacancy.JobVacancyApplies
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
