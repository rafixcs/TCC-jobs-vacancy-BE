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

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	companyDomain := companyfactory.CreateCompanyDomain()

	companiesInfo, err := companyDomain.CompaniesList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&companiesInfo)
}
