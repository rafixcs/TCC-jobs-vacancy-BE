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
	Description string `json:"description"`
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

	err = companyDomain.CreateCompany(createCompanyRequest.Name, createCompanyRequest.Description)
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
