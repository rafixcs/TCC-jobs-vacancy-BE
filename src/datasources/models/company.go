package models

import "time"

type Company struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Email        string    `json:"email"`
	Location     string    `json:"location"`
	CreationDate time.Time `json:"creation_date" db:"creation_date"`
}
