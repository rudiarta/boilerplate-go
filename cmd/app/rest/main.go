package main

import (
	"github.com/joho/godotenv"
	"github.com/rudiarta/boilerplate-go/internal/app/http/rest"
)

func main() {
	err := godotenv.Load("./internal/app/config/.env")
	if err != nil {
		panic(".env is not loaded properly")
	}

	rest.Run()
}
