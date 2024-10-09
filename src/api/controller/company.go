package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/companyfactory"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
)

type CreateCompanyRequest struct {
	Name string `json:"name"`
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

	err = companyDomain.CreateCompany(createCompanyRequest.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type CompanyInfo struct {
	Name         string
	CreationDate time.Time
}

type GetCompaniesResponse struct {
	Companies []CompanyInfo
}

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	companyDomain := companyfactory.CreateCompanyDomain()

	companiesModels, err := companyDomain.CompaniesList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var companiesInfo []CompanyInfo
	for _, companyModel := range companiesModels {
		company := CompanyInfo{
			Name:         companyModel.Name,
			CreationDate: companyModel.CreationDate,
		}

		companiesInfo = append(companiesInfo, company)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&companiesInfo)
}
