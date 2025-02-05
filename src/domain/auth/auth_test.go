package auth

import (
	"testing"

	"github.com/rafixcs/tcc-job-vacancy/src/api/factories/userfactory"
	config "github.com/rafixcs/tcc-job-vacancy/src/configuration"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repoauth"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/domain/company"
	"github.com/stretchr/testify/assert"
)

func setupDB() {
	config.DB_HOST = "127.0.0.1"
	config.DB_USER = "root"
	config.DB_PASSWORD = "root"
	config.DB_NAME = "jobsfinder"
	config.DB_PORT = "5432"
}

func TestAuth(t *testing.T) {
	setupDB()

	userDomain := userfactory.CreateUserDomain()
	mockUser := models.User{
		Name:     "rafael",
		Email:    "test@mail.com",
		Password: "123321",
		Phone:    "123456789",
	}
	_ = userDomain.CreateUser(mockUser.Name, mockUser.Password, mockUser.Email, mockUser.Phone, mockUser.RoleId, company.CompanyInfo{})

	datasource := datasources.DatabasePsql{}
	userRepo := repousers.UserRepository{Datasource: &datasource}
	authRepo := repoauth.AuthRepository{Datasource: &datasource}
	authDomain := AuthDomain{AuthRepo: &authRepo, UserRepo: &userRepo}

	t.Run("Test Login", func(t *testing.T) {
		testEmail := "test@mail.com"
		testPassword := "123321"
		token, _, err := authDomain.UserAuth(testEmail, testPassword)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Wrong user", func(t *testing.T) {
		testEmail := "rafael"
		testPassword := "123321"
		token, _, err := authDomain.UserAuth(testEmail, testPassword)

		assert.NotEmpty(t, err)
		assert.Empty(t, token)
	})

	t.Run("Wrong password", func(t *testing.T) {
		testEmail := "test@mail.com"
		testPassword := "41235123"
		token, _, err := authDomain.UserAuth(testEmail, testPassword)

		assert.NotEmpty(t, err)
		assert.Empty(t, token)
	})

}
