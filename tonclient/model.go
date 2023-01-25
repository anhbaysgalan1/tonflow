package tonclient

import (
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type Wallet struct {
	Address string         `json:"address"`
	Version wallet.Version `json:"version"`
	Seed    string         `json:"seed"`
}
