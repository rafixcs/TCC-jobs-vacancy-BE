package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/controller"
	"github.com/rafixcs/tcc-job-vacancy/src/api/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/rafixcs/tcc-job-vacancy/docs"
)

type JobRouter struct {
	Router *mux.Router
}

func (r *JobRouter) CreateRoutes() {
	r.Router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	r.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"alive"}`))
	}).Methods("GET")
	r.Router.HandleFunc("/api/v1/user", controller.CreateUser).Methods("POST")
	r.Router.HandleFunc("/api/v1/user", controller.GetUserDetails).Methods("GET")
	r.Router.HandleFunc("/api/v1/user", controller.UpdateUser).Methods("PUT")
	r.Router.HandleFunc("/api/v1/user/password", controller.ChangePassword).Methods("PUT")
	r.Router.HandleFunc("/api/v1/auth", controller.Auth).Methods("POST")
	r.Router.HandleFunc("/api/v1/logout", controller.Logout).Methods("POST")
	r.Router.HandleFunc("/api/v1/company", controller.CreateCompany).Methods("POST")
	r.Router.HandleFunc("/api/v1/company/jobs", controller.GetCompanyJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/companies", controller.GetCompanies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job", controller.CreateJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/v1/job/apply", controller.RegisterUserApplyJobVacancy).Methods("POST")
	r.Router.HandleFunc("/api/v1/job/user", controller.GetUserJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/search", controller.SearchJobVacancies).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/applies", controller.GetUsersAppliesToJobVacancy).Methods("GET")
	r.Router.HandleFunc("/api/v1/job/{id}", controller.GetJobVacancyDetails).Methods("GET")
	r.Router.Use(middleware.ContentTypeApplicationJsonMiddleware)
}
