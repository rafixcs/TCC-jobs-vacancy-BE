package models

import "time"

type CompanyModels struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date" db:"creation_date"`
}
