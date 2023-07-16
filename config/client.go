package config

import "github.com/ethereum/go-ethereum/ethclient"

type client struct {
	Network string `json:"network"`
}

func (cli *client) Client() *ethclient.Client {
	client, err := ethclient.Dial(cli.Network)
	if err != nil {
		panic(err)
	}
	return client
}
