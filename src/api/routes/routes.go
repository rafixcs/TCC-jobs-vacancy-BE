package routes

import (
	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/controller"
)

type JobRouter struct {
	Router *mux.Router
}

func (r *JobRouter) CreateRoutes() {
	r.Router.HandleFunc("/api/user", controller.CreateUser).Methods("POST")
}
