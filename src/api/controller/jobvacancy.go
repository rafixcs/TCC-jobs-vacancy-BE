package controller

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
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

	datasource.Open()
	err = datasource.GetError()
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusUnauthorized)
		return
	}
	defer datasource.Close()
	db := datasource.GetDB()

	userApply := models.UserApplies{
		Id:           uuid.NewString(),
		UserId:       userId,
		JobVacancyId: jobId,
	}

	query := `INSERT INTO user_applies (id, job_vacancy_id, user_id) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, userApply.Id, userApply.JobVacancyId, userApply.UserId)
	if err != nil {
		http.Error(w, "execute query failled: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
