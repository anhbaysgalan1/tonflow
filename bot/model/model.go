package model

import (
	"time"
	"tonflow/tonclient"
)

type User struct {
	ID             int64             `json:"id"`
	Username       string            `json:"username"`
	FirstName      string            `json:"firstName"`
	LastName       string            `json:"lastName"`
	LanguageCode   string            `json:"languageCode"`
	Wallet         *tonclient.Wallet `json:"wallet"`
	StageData      *StageData        `json:"stageData"`
	FirstMessageAt time.Time         `json:"firstMessageAt"`
}

type StageData struct {
	Stage         Stage  `json:"stage"`
	AddressToSend string `json:"addressToSend"`
	AmountToSend  string `json:"amountToSend"`
}

type Stage uint8

const (
	ZeroStage Stage = iota
	AddressWait
	AmountWait
	ConfirmationWait
)

func (s Stage) String() string {
	switch s {
	case 0:
		return "zero stage"
	case 1:
		return "address waiting"
	case 2:
		return "amount waiting"
	case 3:
		return "confirmation waiting"
	default:
		return "undefined stage"
	}
}

type UserCache struct {
	UserID int64
	Data   *User
}
