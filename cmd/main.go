package main

import (
	"log"
	"park-wallet/internal/app"
)

func main() {
	tonFlow, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	tonFlow.Run()
}
