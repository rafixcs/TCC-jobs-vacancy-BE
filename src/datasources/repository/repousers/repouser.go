package repousers

import (
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IUserRepository interface {
	Create(user models.User) error
	CheckIfExists(name, email string) (bool, error)
	FindUserByEmail(email string) (models.User, error)
	FindUserById(userId string) (models.User, error)
	UpdateUser(user models.User) error
	UpdatePassword(user models.User) error
}

type UserRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *UserRepository) Create(user models.User) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO users(id, name, email, password, phone, role_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, user.Id, user.Name, user.Email, user.Password, user.Phone, user.RoleId)
	return err
}

func (r *UserRepository) CheckIfExists(name, email string) (bool, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return false, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var alreadyCreatedUser bool
	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM users WHERE name = $1 OR email = $2`
	row := db.QueryRow(query, name, email)
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

func (r *UserRepository) UpdateUser(user models.User) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `UPDATE users SET name = $1, phone = $2 WHERE id = $3`

	_, err = db.Exec(query, user.Name, user.Phone, user.Id)
	return err
}

func (r *UserRepository) UpdatePassword(user models.User) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `UPDATE users SET password = $1 WHERE id = $2`

	_, err = db.Exec(query, user.Password, user.Id)
	return err
}
