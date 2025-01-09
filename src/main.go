package main

import "github.com/rafixcs/tcc-job-vacancy/src/cmd"

// @title Jobs Vacancy Server
// @version 1.0
// @description This is the server for the jobs vacancy project for a post graduation in full stack development
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	cmd.Run()
}
