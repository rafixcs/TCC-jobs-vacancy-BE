package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/repository/repousers"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
)

type IUserDomain interface {
	CreateUser(name, password string, roleId int) error
}

type UserDomain struct {
	UserRepo repousers.IUserRepository
}

func (d *UserDomain) CreateUser(name, password string, roleId int) error {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	alreadyCreatedUser, err := d.UserRepo.CheckIfExists(name)
	if err != nil {
		return err
	}

	if alreadyCreatedUser {
		return fmt.Errorf(`user already created`)
	}

	userId := uuid.NewString()

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	err = d.UserRepo.Create(userId, name, hashedPassword, roleId)

	if err != nil {
		return err
	}

	return nil
}
