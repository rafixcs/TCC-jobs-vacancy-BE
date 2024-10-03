package users

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
	"github.com/rafixcs/tcc-job-vacancy/src/utils"
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
	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM users WHERE name = $1`
	row := db.QueryRow(query, name)
	err = row.Scan(&alreadyCreatedUser)
	if err != nil {
		return err
	}

	if alreadyCreatedUser {
		return fmt.Errorf(`user already created`)
	}

	userId := uuid.NewString()
	roleId := 2 // standard user

	query = `INSERT INTO users(id, name, password, role_id) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(query, userId, name, password, roleId)

	if err != nil {
		return err
	}

	return nil
}

func UserAuth(name, password string) (string, error) {
	err := UserPasswordValidation(name, password)
	if err != nil {
		return "", err
	}

	db, err := datasources.OpenDb()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var (
		id     string
		_name  string
		pass   string
		roleId int
	)
	query := `SELECT id, name, password, role_id  FROM users WHERE name = $1`
	err = db.QueryRow(query, name).Scan(&id, &_name, &pass, &roleId)
	if err != nil {
		return "", err
	}

	if id == "" {
		return "", fmt.Errorf("user not found")
	}

	userModel := models.UserModels{
		Id:       id,
		Name:     _name,
		Password: pass,
		RoleId:   roleId,
	}

	if userModel.Password != password {
		return "", fmt.Errorf("user/password not matching")
	}

	userLoginModel := models.UserLoginsModel{
		Id:        uuid.NewString(),
		UserId:    userModel.Id,
		LoginDate: time.Now(),
	}

	query = `INSERT INTO user_logins(id, user_id, login_date) VALUES ($1, $2, $3);`
	_, err = db.Exec(query, userLoginModel.Id, userLoginModel.UserId, userLoginModel.LoginDate)
	if err != nil {
		return "", err
	}

	token, err := utils.CreateUserJwtToken(userModel.Id, userLoginModel.Id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func Logout(tokenHeader string) error {
	token, err := utils.ParseToken(tokenHeader)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	loginId, ok := claims["login_id"].(string)
	if !ok {
		return fmt.Errorf("token missing login field")
	}

	userId, ok := claims["user_id"].(string)
	if !ok {
		return fmt.Errorf("token missing user field")
	}

	db, err := datasources.OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()

	var validateUserLogin bool
	query := `SELECT CASE WHEN COUNT(*) > 0 THEN true ELSE false END FROM user_logins WHERE id = $1 AND user_id = $2`
	err = db.QueryRow(query, loginId, userId).Scan(&validateUserLogin)
	if err != nil {
		return err
	}

	if !validateUserLogin {
		return fmt.Errorf("user login not found")
	}

	query = `UPDATE user_logins SET logout_date=$1 WHERE id = $2 AND user_id = $3`
	_, err = db.Exec(query, time.Now(), loginId, userId)
	if err != nil {
		return err
	}

	return nil
}
