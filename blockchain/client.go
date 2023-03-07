package blockchain

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type Client struct {
	liteClient *liteclient.ConnectionPool
	tonClient  *ton.APIClient
}

func NewClient(configUrl string, prod bool) (*Client, error) {
	var (
		liteClient *liteclient.ConnectionPool
		err        error
	)

	switch prod {
	case true:
		liteClient = liteclient.NewConnectionPool()
		err = liteClient.AddConnection(context.Background(), "116.203.233.170:11358", "VdZyqnuUGqO9BaF2v+lt7isk/igihPUu9Vh74/wuwrc=")
	default:
		liteClient = liteclient.NewConnectionPool()
		err = liteClient.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	}
	if err != nil {
		return nil, err
	}

	return &Client{
		liteClient: liteClient,
		tonClient:  ton.NewAPIClient(liteClient),
	}, nil
}
