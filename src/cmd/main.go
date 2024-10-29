package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rafixcs/tcc-job-vacancy/src/api/middleware"
	"github.com/rafixcs/tcc-job-vacancy/src/api/routes"
	config "github.com/rafixcs/tcc-job-vacancy/src/configuration"
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
	r.Use(middleware.LoggingMiddleware)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	myRouter := routes.JobRouter{Router: r}
	myRouter.CreateRoutes()

	log.Printf(`Server listening on port %v`, config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)))
}
