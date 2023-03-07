package blockchain

import (
	"context"
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"strconv"
	"strings"
	"tonflow/model"
	"tonflow/pkg"
)

func (c *Client) NewWallet() (*model.Wallet, error) {
	seed := wallet.NewSeed()
	version := wallet.V4R2

	w, err := wallet.FromSeed(c.tonClient, seed, version)
	if err != nil {
		return nil, err
	}

	return &model.Wallet{
		Address: w.Address().String(),
		Version: version,
		Seed:    strings.Join(seed, " "),
	}, nil
}

func (c *Client) GetWalletBalance(ctx context.Context, addr string) (string, error) {
	ctx = c.liteClient.StickyContext(ctx)

	wltAddr, err := address.ParseAddr(addr)
	if err != nil {
		return "", err
	}

	block, err := c.tonClient.GetMasterchainInfo(ctx)
	if err != nil {
		return "", err
	}

	wlt, err := c.tonClient.GetAccount(ctx, block, wltAddr)
	if err != nil {
		return "", err
	}

	if !wlt.IsActive {
		return "0", nil
	}

	return wlt.State.Balance.TON(), nil
}

func (c *Client) ValidateWallet(addr string) error {
	_, err := address.ParseAddr(addr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Send(ctx context.Context, user *model.User) error {
	seed, err := pkg.Decode(user.Wallet.Seed, strconv.FormatInt(user.ID, 10))
	if err != nil {
		return err
	}

	words := strings.Split(seed, " ")
	w, err := wallet.FromSeed(c.tonClient, words, user.Wallet.Version)
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

func (c *Client) SendAll(ctx context.Context, user *model.User) error {
	seed, err := pkg.Decode(user.Wallet.Seed, strconv.FormatInt(user.ID, 10))
	if err != nil {
		return err
	}

	words := strings.Split(seed, " ")
	w, err := wallet.FromSeed(c.tonClient, words, user.Wallet.Version)
	if err != nil {
		return err
	}

	var body *cell.Cell
	if user.StageData.Comment != "" {
		body, err = CreateCommentCell(user.StageData.Comment)
		if err != nil {
			return err
		}
	}

	err = w.Send(ctx, &wallet.Message{
		Mode: 128,
		InternalMessage: &tlb.InternalMessage{
			IHRDisabled: true,
			Bounce:      false,
			DstAddr:     address.MustParseAddr(user.StageData.AddressToSend),
			Amount:      tlb.MustFromTON("0"),
			Body:        body,
		},
	}, true)
	if err != nil {
		return err
	}

	return nil
}

func CreateCommentCell(text string) (*cell.Cell, error) {
	// comment ident
	root := cell.BeginCell().MustStoreUInt(0, 32)

	if err := root.StoreStringSnake(text); err != nil {
		return nil, fmt.Errorf("failed to build comment: %w", err)
	}

	return root.EndCell(), nil
}
