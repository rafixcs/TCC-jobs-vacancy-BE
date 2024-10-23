package repousers

import (
	"fmt"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IUserRepository interface {
	Create(userId, name, password, email string, roleId int) error
	CheckIfExists(name string) (bool, error)
	FindUser(name string) (models.UserModels, error)
}

type UserRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *UserRepository) Create(userId, name, password, email string, roleId int) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO users(id, name, email, password, role_id) VALUES ($1, $2, $3, $4, $5)`
	_, err = db.Exec(query, userId, name, email, password, roleId)
	return err
}

func (r *UserRepository) CheckIfExists(name string) (bool, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return false, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var alreadyCreatedUser bool
	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM users WHERE name = $1`
	row := db.QueryRow(query, name)
	err = row.Scan(&alreadyCreatedUser)
	if err != nil {
		return false, err
	}

	return alreadyCreatedUser, nil
}

func (r *UserRepository) FindUser(name string) (models.UserModels, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.UserModels{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var (
		id     string
		pass   string
		roleId int
	)
	query := `SELECT id, password, role_id  FROM users WHERE name = $1`
	err = db.QueryRow(query, name).Scan(&id, &pass, &roleId)
	if err != nil {
		return models.UserModels{}, err
	}

	if id == "" {
		return models.UserModels{}, fmt.Errorf("user not found")
	}

	userModel := models.UserModels{
		Id:       id,
		Name:     name,
		Password: pass,
		RoleId:   roleId,
	}

	return userModel, nil
}
