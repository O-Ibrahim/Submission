package main

import (
	"log"
	"takehome/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Panic(err)
	}
	application.Run()
}
