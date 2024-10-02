package users

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
)

type User struct {
	Name     string
	Password string
}

func UserPasswordValidation(name, password string) error {
	if name == "" || len(name) < 3 {
		return fmt.Errorf("invalid user name provided")
	}

	if password == "" || len(password) < 6 {
		return fmt.Errorf("invalid user password")
	}

	return nil
}

func CreateUser(name, password string) error {

	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	db, err := datasources.OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()

	var alreadyCreatedUser bool
	query := `SELECT IF(COUNT(*), true, false) FROM users WHERE name=?`
	err = db.QueryRow(query, name).Scan(&alreadyCreatedUser)
	if err != nil {
		return err
	}

	if alreadyCreatedUser {
		return fmt.Errorf(`user already created`)
	}

	userId := uuid.NewString()
	roleId := 2 // standard user

	query = `INSERT INTO users(id, name, password, role_id) VALUES ($1, $2, $3, $4);`
	_, err = db.Exec(query, userId, name, password, roleId)

	if err != nil {
		return err
	}

	return nil
}

func UserAuth(name, password string) error {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return err
	}

	return nil
}
