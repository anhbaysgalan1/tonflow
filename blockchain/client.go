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

	//lt := liteclient.NewConnectionPool()
	//err = lt.AddConnection(context.Background(), "ip:port", "key")
	//if err != nil {
	//	log.Fatalf("add connection error: %s", err)
	//}

	tonClient := ton.NewAPIClient(liteClient)

	return &Client{
		liteClient: liteClient,
		tonClient:  tonClient,
	}, nil
}
