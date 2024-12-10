package cmd

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

// @title Jobs Vacancy Server
// @version 1.0
// @description This is the server for the jobs vacancy project for a post graduation in full stack development
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func Run() {

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

	config.PORT = "8080"
	log.Printf(`Server listening on port %v`, config.PORT)
	log.Fatal(http.ListenAndServe(":"+config.PORT, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)))
}
