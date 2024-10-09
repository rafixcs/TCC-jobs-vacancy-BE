package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type CreateJobVacancyRequest struct {
	CompanyName string `json:"company"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

func CreateJobVacancy(w http.ResponseWriter, r *http.Request) {
	tokenHeader := r.Header.Get("Authorization")
	userId, err := getUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}

	var requestContent CreateJobVacancyRequest
	json.NewDecoder(r.Body).Decode(&requestContent)

	datasource := datasources.DatabasePsql{}
	repoCompany := repocompany.CompanyRepository{Datasource: &datasource}
	company, err := repoCompany.FindCompanyByName(requestContent.CompanyName)
	if err != nil {
		http.Error(w, "failed to get company", http.StatusInternalServerError)
		return
	}

	jobVacancy := models.JobVacancy{
		Id:           uuid.NewString(),
		UserId:       userId,
		CompanyId:    company.Id,
		Title:        requestContent.Title,
		Description:  requestContent.Description,
		CreationDate: time.Now(),
	}

	query := "INSERT INTO job_vacancies(id, user_id, company_id, description, title, creation_date) VALUES ($1, $2, $3, $4, $5, $6)"

	datasource.Open()
	err = datasource.GetError()
	if err != nil {
		http.Error(w, "failed to open", http.StatusInternalServerError)
		return
	}
	defer datasource.Close()
	db := datasource.GetDB()

	_, err = db.Exec(query, jobVacancy.Id, jobVacancy.UserId, jobVacancy.CompanyId, jobVacancy.Description, jobVacancy.Title, jobVacancy.CreationDate)
	if err != nil {
		http.Error(w, "failed to insert: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getUserIdFromToken(tokenHeader string) (string, error) {
	token, err := utils.ParseToken(tokenHeader)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("token missing user field")
	}

	return userId, nil
}

func RegisterUserApplyJobVacancy(w http.ResponseWriter, r *http.Request) {
	/*tokenHeader := r.Header.Get("Authorization")
	userId, err := getUserIdFromToken(tokenHeader)
	if err != nil {
		http.Error(w, "failed to parse Authorization token", http.StatusUnauthorized)
		return
	}*/
}
