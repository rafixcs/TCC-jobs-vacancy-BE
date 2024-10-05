package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/routes"
	"github.com/rafixcs/tcc-job-vacancy/src/datasources"
)

func main() {

	database := datasources.DatabasePsql{}
	database.Open()
	err := database.GetError()

	if err != nil {
		panic(err)
	}
	database.Close()

	r := mux.NewRouter()

	myRouter := routes.JobRouter{Router: r}
	myRouter.CreateRoutes()

	log.Println("Server listening on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
