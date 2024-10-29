package config

import "os"

var (
	DB_HOST      = os.Getenv("DB_HOST")
	DB_USER      = os.Getenv("DB_USER")
	DB_PASSWORD  = os.Getenv("DB_PASSWORD")
	DB_NAME      = os.Getenv("DB_NAME")
	DB_PORT      = os.Getenv("DB_PORT")
	TOKEN_SECRET = os.Getenv("TOKEN_SECRET")
	PORT         = os.Getenv("PORT")

	CF_ACCESS_KEY        = os.Getenv("CF_ACCESS_KEY")
	CF_SECRET_ACCESS_KEY = os.Getenv("CF_SECRET_ACCESS_KEY")
	R2_ENDPOINT          = os.Getenv("R2_ENDPOINT")
	R2_BUCKET            = os.Getenv("BUCKET_NAME")
)
