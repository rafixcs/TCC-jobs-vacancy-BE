package users

import "fmt"

type User struct {
	Name     string
	Password string
}

func CreateUser(name, password string) error {

	if name == "" || len(name) < 3 {
		return fmt.Errorf("Invalid user name provided")
	}

	if password == "" || len(password) < 6 {
		return fmt.Errorf("Invalid user password")
	}

	return nil
}
