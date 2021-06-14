package context

import "os"

type ctx struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	APP_SALT    string
}

var GlobalCtx = ctx{
	DB_HOST:     os.Getenv("DB_HOST"),
	DB_PORT:     os.Getenv("DB_PORT"),
	DB_USER:     os.Getenv("DB_USER"),
	DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	DB_NAME:     os.Getenv("DB_NAME"),
	APP_SALT:    os.Getenv("APP_SALT"),
}
