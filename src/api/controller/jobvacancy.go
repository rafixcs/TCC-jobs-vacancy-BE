package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/jobfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/jobvacancy"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type CreateJobVacancyRequest struct {
	CompanyName string `json:"company"`
	Description string `json:"description"`
	Title       string `json:"title"`
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

	err = jobvacancyDomain.CreateJobVacancy(userId, requestContent.CompanyName, requestContent.Description, requestContent.Title)
	if err != nil {
		http.Error(w, "failed to create job vacancy", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

	jobvacancyDomain := jobfactory.CreateJobVacancyDomain()

	err = jobvacancyDomain.CreateUserJobApply(userId, jobId)
	if err != nil {
		http.Error(w, "could not create user apply", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type GetCompaniesJobVacanciesResponse struct {
	CompanyName  string                  `json:"company"`
	JobVacancies []jobvacancy.JobVacancy `json:"job_vacancies"`
}

func GetCompanyJobVacancies(w http.ResponseWriter, r *http.Request) {
	companyName := r.URL.Query().Get("company")
	if companyName == "" {
		http.Error(w, "missing company name", http.StatusBadRequest)
		return
	}

	jobvacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobVacancies, err := jobvacancyDomain.GetCompanyJobVacancies(companyName)
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
	JobApplies []jobvacancy.JobVacancy
}

func GetUserJobVacancies(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := utils.GetUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	jobvacancyDomain := jobfactory.CreateJobVacancyDomain()
	jobApplies, err := jobvacancyDomain.GetUserJobApplies(userId)
	if err != nil {
		http.Error(w, "failed to get user job applies list", http.StatusBadRequest)
		return
	}

	responseBody := UserJobAppliesResponse{
		JobApplies: jobApplies,
	}

	json.NewEncoder(w).Encode(&responseBody)
}
