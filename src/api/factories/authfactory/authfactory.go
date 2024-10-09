package authfactory

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repoauth"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/auth"
)

func CreateAuthDomain() auth.IAuthDomain {
	datasource := datasources.DatabasePsql{}
	authRepo := repoauth.AuthRepository{Datasource: &datasource}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	authDomain := auth.AuthDomain{AuthRepo: &authRepo, UserRepo: &userRepo}

	return &authDomain
}
