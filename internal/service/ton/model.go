package ton

import (
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type Transaction struct {
	Source  string
	Hash    string
	Value   int
	Comment string
}

type Address struct {
	Address    string
	Valid      bool
	Status     *tlb.AccountStatus
	LastTxLT   uint64
	LastTxHash []byte
	Balance    string
}

type Wallet struct {
	Address string
	Version wallet.Version
	Seed    string
}
