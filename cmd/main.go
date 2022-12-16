package main

import (
	"log"
	"ton-flow-bot/internal/app"
)

func main() {
	tonFlow, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	tonFlow.Run()
}
