package users

import "fmt"

func UserPasswordValidation(name, password string) error {
	if name == "" || len(name) < 3 {
		return fmt.Errorf("invalid user name provided")
	}

	if password == "" || len(password) < 6 {
		return fmt.Errorf("invalid user password")
	}

	return nil
}
