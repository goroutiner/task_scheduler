package config

import "os"

var (
	Port     = os.Getenv("PORT")
	Mode     = os.Getenv("MODE")
	PsqlUrl  = os.Getenv("DATABASE_URL")
	Password = os.Getenv("PASSWORD")
)
