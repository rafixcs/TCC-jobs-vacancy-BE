package jobfactory

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repojobvacancy"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/jobvacancy"
)

func CreateJobVacancyDomain() jobvacancy.IJobVacancyDomain {
	datasource := datasources.DatabasePsql{}
	repoCompany := repocompany.CompanyRepository{Datasource: &datasource}
	jobvacancyRepo := repojobvacancy.JobVacancyRepository{Datasource: &datasource}
	jobvacancyDomain := jobvacancy.JobVacancyDomain{JobVacancyRepo: &jobvacancyRepo, CompanyRepo: &repoCompany}
	return jobvacancyDomain
}
