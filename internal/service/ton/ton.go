package ton

import (
	"context"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

type Ton struct {
	liteClient *liteclient.ConnectionPool
	tonAPI     *ton.APIClient
}

func NewTon(configUrl string) (*Ton, error) {
	liteClient := liteclient.NewConnectionPool()
	err := liteClient.AddConnectionsFromConfigUrl(context.Background(), configUrl)
	if err != nil {
		return nil, err
	}

	tonAPIClient := ton.NewAPIClient(liteClient)

	return &Ton{
		liteClient: liteClient,
		tonAPI:     tonAPIClient,
	}, nil
}
