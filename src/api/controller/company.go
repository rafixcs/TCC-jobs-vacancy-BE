package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
)

type CreateCompanyRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

// CreateCompany godoc
// @Summary Create company
// @Description Create company
// @Tags Company
// @Param Authorization header string true "Authorization token"
// @Param createcompanyrequest body CreateCompanyRequest true "Create company"
// @Success 201 "Created company"
// @Failure 400 "Bad request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal server error"
// @Router /api/v1/company [post]
func CreateCompany(w http.ResponseWriter, r *http.Request) {
	var createCompanyRequest CreateCompanyRequest

	err := json.NewDecoder(r.Body).Decode(&createCompanyRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	datasource := datasources.DatabasePsql{}
	companyRepo := repocompany.CompanyRepository{Datasource: &datasource}
	companyDomain := company.CompanyDomain{CompanyRepo: &companyRepo}

	_, err = companyDomain.CreateCompany(createCompanyRequest.Name, createCompanyRequest.Email, createCompanyRequest.Description, createCompanyRequest.Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type GetCompaniesResponse struct {
	Companies []company.CompanyInfo
}

// CreateCompany godoc
// @Summary Create company
// @Description Create company
// @Tags Company
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} GetCompaniesResponse "Success"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /api/v1/companies [get]
func GetCompanies(w http.ResponseWriter, r *http.Request) {
	companyDomain := companyfactory.CreateCompanyDomain()

	companiesInfo, err := companyDomain.CompaniesList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&companiesInfo)
}
