package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"
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

	datasource := datasources.DatabasePsql{}
	repoCompany := repocompany.CompanyRepository{Datasource: &datasource}
	jobvacancyRepo := repojobvacancy.JobVacancyRepository{Datasource: &datasource}
	jobvacancyDomain := jobvacancy.JobVacancyDomain{JobVacancyRepo: &jobvacancyRepo, CompanyRepo: &repoCompany}

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

	datasource := datasources.DatabasePsql{}
	repoCompany := repocompany.CompanyRepository{Datasource: &datasource}
	jobvacancyRepo := repojobvacancy.JobVacancyRepository{Datasource: &datasource}
	jobvacancyDomain := jobvacancy.JobVacancyDomain{JobVacancyRepo: &jobvacancyRepo, CompanyRepo: &repoCompany}

	err = jobvacancyDomain.CreateUserJobApply(userId, jobId)
	if err != nil {
		http.Error(w, "could not create user apply", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
