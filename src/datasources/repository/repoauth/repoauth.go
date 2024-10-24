package repoauth

import (
	"time"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type IAuthRepository interface {
	CreateLogin(userLoginModel models.UserLogins) error
	ValidateLogin(loginId, userId string) (bool, error)
	UpdateToLogout(loginId, userId string) error
}

type AuthRepository struct {
	Datasource datasources.IDatabasePsql
}

func (r *AuthRepository) CreateLogin(userLoginModel models.UserLogins) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `INSERT INTO user_logins(id, user_id, login_date) VALUES ($1, $2, $3);`
	_, err = db.Exec(query, userLoginModel.Id, userLoginModel.UserId, userLoginModel.LoginDate)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) ValidateLogin(loginId, userId string) (bool, error) {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return false, err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	var validateUserLogin bool
	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM user_logins WHERE id = $1 AND user_id = $2`
	err = db.QueryRow(query, loginId, userId).Scan(&validateUserLogin)
	if err != nil {
		return false, err
	}

	return validateUserLogin, nil
}

func (r *AuthRepository) UpdateToLogout(loginId, userId string) error {
	r.Datasource.Open()
	err := r.Datasource.GetError()
	if err != nil {
		return err
	}
	defer r.Datasource.Close()
	db := r.Datasource.GetDB()

	query := `UPDATE user_logins SET logout_date=$1 WHERE id = $2 AND user_id = $3`
	_, err = db.Exec(query, time.Now(), loginId, userId)
	if err != nil {
		return err
	}

	return nil
}
