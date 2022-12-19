package ton

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
)

func (ton *Ton) NewWallet() (*Wallet, error) {
	seed := wallet.NewSeed()

	version := wallet.V4R2

	w, err := wallet.FromSeed(ton.tonAPI, seed, version)
	if err != nil {
		return nil, err
	}

	seedString := strings.Join(seed, " ")
	return &Wallet{
		Address: w.Address().String(),
		Version: version,
		Seed:    seedString,
	}, nil
}

func (ton *Ton) GetWalletBalance(wallet string) (string, error) {
	ctx := ton.liteClient.StickyContext(context.Background())

	addr, err := address.ParseAddr(wallet)
	if err != nil {
		log.Error().Err(err).Msg("parse TON address")
		return "", err
	}

	block, err := ton.tonAPI.GetMasterchainInfo(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get block")
		return "", err
	}

	wlt, err := ton.tonAPI.GetAccount(ctx, block, addr)
	if err != nil {
		log.Error().Err(err).Msg("get account info")
		return "", err
	}

	if wlt.IsActive {
		return wlt.State.Balance.TON(), nil
	}

	return "0", nil
}
