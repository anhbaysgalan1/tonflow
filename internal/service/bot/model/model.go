package model

import "time"

type User struct {
	ID             int64
	Username       string
	FirstName      string
	LastName       string
	LanguageCode   string
	Wallet         string
	FirstMessageAt time.Time
	LastMessageAt  time.Time
}
