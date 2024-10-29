package jobvacancy

import (
	"time"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type JobVacancyInfo struct {
	Id           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	Company      string    `json:"company"`
	Location     string    `json:"location"`
	UrlResume    string    `json:"url_resume"`
}

func (j JobVacancyInfo) TrasformFromModel(model models.JobVacancy, companyName string) JobVacancyInfo {

	info := JobVacancyInfo{
		Id:           model.Id,
		Title:        model.Title,
		Description:  model.Description,
		CreationDate: model.CreationDate,
		Company:      companyName,
		Location:     model.Location,
	}

	return info
}

func (j JobVacancyInfo) TransformSliceModelCompany(model []models.JobVacancy, companiesNames []string) []JobVacancyInfo {
	var jobVacanciesInfo []JobVacancyInfo

	for i, job := range model {
		companyName := companiesNames[i]
		info := j.TrasformFromModel(job, companyName)
		jobVacanciesInfo = append(jobVacanciesInfo, info)
	}

	return jobVacanciesInfo
}

func (j JobVacancyInfo) TransformSliceModel(model []models.JobVacancy) []JobVacancyInfo {
	var jobVacanciesInfo []JobVacancyInfo

	for _, job := range model {
		info := j.TrasformFromModel(job, "")
		jobVacanciesInfo = append(jobVacanciesInfo, info)
	}

	return jobVacanciesInfo
}
