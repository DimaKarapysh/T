package app

import (
	"T/bootstrap"
	"github.com/joho/godotenv"
)

func InitApp() (err error) {
	err = godotenv.Load()
	if err != nil {
		return
	}

	err = bootstrap.InitRest()
	if err != nil {
		return
	}
	err = nil
	return
}
