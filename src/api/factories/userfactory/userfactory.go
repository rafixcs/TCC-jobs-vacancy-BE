package userfactory

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repocompany"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/users"
)

func CreateUserDomain() users.IUserDomain {
	datasource := datasources.DatabasePsql{}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	companyRepo := repocompany.CompanyRepository{Datasource: &datasource}
	userDomain := users.UserDomain{UserRepo: &userRepo, CompanyRepo: &companyRepo}
	return &userDomain
}
