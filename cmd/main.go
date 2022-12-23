package main

import (
	"log"
	"tonflow/internal/app"
)

func main() {
	flowWallet, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	flowWallet.Run()
}
