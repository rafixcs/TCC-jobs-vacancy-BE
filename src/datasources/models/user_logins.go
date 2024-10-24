package models

import "time"

type UserLogins struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	LoginDate  time.Time `json:"login_date"`
	LogoutDate time.Time `json:"logout_date"`
}
