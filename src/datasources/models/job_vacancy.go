package models

import "time"

type JobVacancy struct {
	Id               string
	UserId           string    `json:"user_id"`
	CompanyId        string    `json:"company_id"`
	Description      string    `json:"description"`
	Title            string    `json:"title"`
	Location         string    `json:"location"`
	CreationDate     time.Time `json:"creation_date"`
	Salary           string    `json:"salary"`
	Requirements     string    `json:"requirements"`
	Responsabilities string    `json:"responsabilities"`
}

type UserApplies struct {
	Id           string
	UserId       string `json:"user_id"`
	JobVacancyId string `json:"job_vacancy_id"`
}
