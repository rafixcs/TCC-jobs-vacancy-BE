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
	r.Router.HandleFunc("/api/user", controller.CreateUser).Methods("POST")
	r.Router.HandleFunc("/api/auth", controller.Auth).Methods("POST")
	r.Router.HandleFunc("/api/logout", controller.Logout).Methods("POST")
	r.Router.HandleFunc("/api/company", controller.CreateCompany).Methods("POST")
	r.Router.HandleFunc("/api/companies", controller.GetCompanies).Methods("GET")
	r.Router.HandleFunc("/api/job", controller.CreateJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/job/apply", controller.RegisterUserApplyJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/job/company", controller.GetCompanyJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/job/user", controller.GetUserJobVacancies).Methods("GET")
	r.Router.Use(middleware.ContentTypeApplicationJsonMiddleware)
}
