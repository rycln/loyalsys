package main

import (
	"log"

	"github.com/rycln/loyalsys/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
