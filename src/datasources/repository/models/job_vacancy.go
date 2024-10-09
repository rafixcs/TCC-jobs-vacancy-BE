package models

import "time"

type JobVacancy struct {
	Id           string
	UserId       string    `json:"user_id"`
	CompanyId    string    `json:"company_id"`
	Description  string    `json:"description"`
	Title        string    `json:"title"`
	CreationDate time.Time `json:"creation_date"`
}
