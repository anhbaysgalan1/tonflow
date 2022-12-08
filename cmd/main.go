package main

import (
	"log"
	"ton-flow-bot/internal/app"
)

func main() {
	botApp, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	botApp.Run()
}
