package company

import (
	"time"

	"github.com/rafixcs/tcc-job-vacancy/src/datasources/models"
)

type CompanyInfo struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	CreationDate time.Time `json:"creation_date"`
	Description  string    `json:"description"`
	Email        string    `json:"email"`
	Location     string    `json:"location"`
}

func (CompanyInfo) TransformFromModel(companyModel models.Company) CompanyInfo {
	return CompanyInfo{
		Id:           companyModel.Id,
		Name:         companyModel.Name,
		CreationDate: companyModel.CreationDate,
		Description:  companyModel.Description,
		Email:        companyModel.Email,
		Location:     companyModel.Location,
	}
}
