package repousers

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IUserRepository interface {
	Create(userId, name, password, email, phone string, roleId int) error
	CheckIfExists(name string) (bool, error)
	FindUserByEmail(email string) (models.User, error)
	FindUserById(userId string) (models.User, error)
}

type UserRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *UserRepository) Create(userId, name, password, email, phone string, roleId int) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO users(id, name, email, password, phone, role_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, userId, name, email, password, phone, roleId)
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

func (r *UserRepository) FindUserByEmail(email string) (models.User, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.User{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var userModel models.User
	query := `SELECT id, name, password, role_id, email, phone  FROM users WHERE email = $1`
	err = db.QueryRow(query, email).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Password,
		&userModel.RoleId,
		&userModel.Email,
		&userModel.Phone,
	)

	if err != nil {
		return models.User{}, err
	}

	return userModel, nil
}

func (r *UserRepository) FindUserById(userId string) (models.User, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return models.User{}, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var userModel models.User
	query := `SELECT id, name, password, email, phone, role_id  FROM users WHERE id = $1`
	err = db.QueryRow(query, userId).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Password,
		&userModel.Email,
		&userModel.Phone,
		&userModel.RoleId,
	)

	if err != nil {
		return models.User{}, err
	}

	return userModel, nil
}
