package model

import (
	"github.com/xssnick/tonutils-go/ton/wallet"
	"time"
)

type User struct {
	ID             int64      `json:"id"`
	Username       string     `json:"username"`
	FirstName      string     `json:"firstName"`
	LastName       string     `json:"lastName"`
	LanguageCode   string     `json:"languageCode"`
	Wallet         *Wallet    `json:"wallet"`
	StageData      *StageData `json:"stageData"`
	IsExisted      bool       `json:"isExisted"`
	FirstMessageAt time.Time  `json:"firstMessageAt"`
}

type Wallet struct {
	Address string         `json:"address"`
	Version wallet.Version `json:"version"`
	Seed    string         `json:"seed"`
}

type StageData struct {
	Stage         Stage  `json:"stage"`
	AddressToSend string `json:"addressToSend"`
	AmountToSend  string `json:"amountToSend"`
	SendAll       bool   `json:"sendAll"`
	Comment       string `json:"comment"`
}

type Stage uint8

const (
	ZeroStage Stage = iota
	AddressWait
	AmountWait
	CommentWait
	ConfirmationWait
)

var EmptyStageData = StageData{
	Stage:         ZeroStage,
	AddressToSend: "",
	AmountToSend:  "",
	SendAll:       false,
	Comment:       "",
}

func (s Stage) String() string {
	switch s {
	case 0:
		return "zero stage"
	case 1:
		return "address waiting"
	case 2:
		return "amount waiting"
	case 3:
		return "comment waiting"
	case 4:
		return "confirmation waiting"
	default:
		return "undefined stage"
	}
}
