package main

import (
	"flow-wallet/internal/app"
	"log"
)

func main() {
	tonFlow, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	tonFlow.Run()
}
