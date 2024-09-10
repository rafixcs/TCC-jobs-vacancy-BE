package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/routes"
)

func main() {

	r := mux.NewRouter()

	myRouter := routes.JobRouter{Router: r}
	myRouter.CreateRoutes()

	log.Fatal(http.ListenAndServe(":8000", r))

}
