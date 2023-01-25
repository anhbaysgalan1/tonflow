package tonclient

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type TonClient struct {
	liteClient *liteclient.ConnectionPool
	tonAPI     *ton.APIClient
}

func NewTonClient(configUrl string) (*TonClient, error) {
	liteClient := liteclient.NewConnectionPool()
	err := liteClient.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return nil, err
	}

	tonAPIClient := ton.NewAPIClient(liteClient)

	return &TonClient{
		liteClient: liteClient,
		tonAPI:     tonAPIClient,
	}, nil
}
