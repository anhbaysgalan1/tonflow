package tonclient

import (
	"context"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"strings"
	"tonflow/model"
)

func (ton *TonClient) NewWallet() (*model.Wallet, error) {
	seed := wallet.NewSeed()
	version := wallet.V4R2

	w, err := wallet.FromSeed(ton.tonAPI, seed, version)
	if err != nil {
		return nil, err
	}

	return &model.Wallet{
		Address: w.Address().String(),
		Version: version,
		Seed:    strings.Join(seed, " "),
	}, nil
}

func (ton *TonClient) ValidateWallet(addr string) error {
	_, err := address.ParseAddr(addr)
	if err != nil {
		return err
	}
	return nil
}

func (ton *TonClient) GetWalletBalance(addr string) (string, error) {
	ctx := ton.liteClient.StickyContext(context.Background())

	wltAddr, err := address.ParseAddr(addr)
	if err != nil {
		return "", err
	}

	block, err := ton.tonAPI.GetMasterchainInfo(ctx)
	if err != nil {
		return "", err
	}

	wlt, err := ton.tonAPI.GetAccount(ctx, block, wltAddr)
	if err != nil {
		return "", err
	}

	if !wlt.IsActive {
		return "0", nil
	}

	return wlt.State.Balance.TON(), nil
}
