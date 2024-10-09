package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/jobfactory"
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
