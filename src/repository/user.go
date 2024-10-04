package repository

import "github.com/rafixcs/tcc-job-vacancy/src/repository/models"

type UserRepository struct {
	RegisteredUsers []models.UserModels
}

func (uc *UserRepository) CreateUser(username, password string) {
	user := models.UserModels{
		Username: username,
		Password: password,
	}

	uc.RegisteredUsers = append(uc.RegisteredUsers, user)
}
