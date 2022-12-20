package main

import (
	"log"
	"park-wallet/internal/app"
)

func main() {
	flowWallet, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	flowWallet.Run()
}
