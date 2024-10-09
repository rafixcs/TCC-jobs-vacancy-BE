package companyfactory

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
)

func CreateCompanyDomain() company.ICompanyDomain {
	datasource := datasources.DatabasePsql{}
	companyRepo := repocompany.CompanyRepository{Datasource: &datasource}
	companyDomain := company.CompanyDomain{CompanyRepo: &companyRepo}
	return &companyDomain
}
