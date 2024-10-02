package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/routes"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
)

func main() {

	db, err := datasources.OpenDb()
	if err != nil {
		panic(err)
	}
	db.Close()

	r := mux.NewRouter()

	myRouter := routes.JobRouter{Router: r}
	myRouter.CreateRoutes()

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
