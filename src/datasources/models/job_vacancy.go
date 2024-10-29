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
	Responsibilities string    `json:"responsibilities"`
	JobType          string    `json:"job_type"`
	ExperienceLevel  string    `json:"experience_level"`
}

type UserApplies struct {
	Id           string
	UserId       string `json:"user_id"`
	JobVacancyId string `json:"job_vacancy_id"`
	CoverLetter  string `json:"cover_letter"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
	UrlResume    string `json:"url_resume"`
}
