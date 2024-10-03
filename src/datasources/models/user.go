package models

type UserModels struct {
	Id       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password"`
	RoleId   int    `json:"role_id" db:"role_id"`
}
