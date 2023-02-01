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

func NewClient(configUrl string) (*Client, error) {
	liteClient := liteclient.NewConnectionPool()
	err := liteClient.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return nil, err
	}

	tonClient := ton.NewAPIClient(liteClient)

	return &Client{
		liteClient: liteClient,
		tonClient:  tonClient,
	}, nil
}
