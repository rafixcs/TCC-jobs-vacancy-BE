package routes

import (
	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/controller"
	"github.com/rafixcs/tcc-job-vacancy/src/api/middleware"
)

type JobRouter struct {
	Router *mux.Router
}

func (r *JobRouter) CreateRoutes() {
	r.Router.HandleFunc("/api/v1/user", controller.CreateUser).Methods("POST")
	r.Router.HandleFunc("/api/v1/auth", controller.Auth).Methods("POST")
	r.Router.HandleFunc("/api/v1/logout", controller.Logout).Methods("POST")
	r.Router.HandleFunc("/api/v1/company", controller.CreateCompany).Methods("POST")
	r.Router.HandleFunc("/api/v1/companies", controller.GetCompanies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job", controller.CreateJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/v1/job/{id}", controller.GetJobVacancyDetails).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/apply", controller.RegisterUserApplyJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/v1/job/company", controller.GetCompanyJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/user", controller.GetUserJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/search", controller.SearchJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/applies", controller.GetUsersAppliesToJobVacancy).Methods("GET")
	r.Router.Use(middleware.ContentTypeApplicationJsonMiddleware)
}
