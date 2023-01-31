package tonclient

import (
	"context"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
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

func (ton *TonClient) GetWalletBalance(ctx context.Context, addr string) (string, error) {
	ctx = ton.liteClient.StickyContext(ctx)

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

func (ton *TonClient) Send(ctx context.Context, user *model.User) error {
	words := strings.Split(user.Wallet.Seed, " ")
	w, err := wallet.FromSeed(ton.tonAPI, words, user.Wallet.Version)
	if err != nil {
		return err
	}

	addr := address.MustParseAddr(user.StageData.AddressToSend)
	amount := tlb.MustFromTON(user.StageData.AmountToSend)
	comment := user.StageData.Comment

	err = w.TransferNoBounce(ctx, addr, amount, comment, true)
	if err != nil {
		return err
	}

	return nil
}
