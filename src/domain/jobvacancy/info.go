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
}

func (j JobVacancyInfo) TrasformFromModel(model models.JobVacancy) JobVacancyInfo {

	info := JobVacancyInfo{
		Id:           model.Id,
		Title:        model.Title,
		Description:  model.Description,
		CreationDate: model.CreationDate,
	}

	return info
}

func (j JobVacancyInfo) TransformSliceModel(model []models.JobVacancy) []JobVacancyInfo {
	var jobVacanciesInfo []JobVacancyInfo

	for _, job := range model {
		info := j.TrasformFromModel(job)
		jobVacanciesInfo = append(jobVacanciesInfo, info)
	}

	return jobVacanciesInfo
}
