package main

import (
	"log"
	"tonflow/internal/app"
)

func main() {
	tonflow, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	tonflow.Run()
}
