package config

import "github.com/ethereum/go-ethereum/ethclient"

type Client struct {
	HttpCli *ethclient.Client
	WsCli   *ethclient.Client
}

type client struct {
	NetworkWS   string `json:"networkWS"`
	NetworkHTTP string `json:"networkHTTP"`
}

func (cli *client) Client() *Client {
	wsClient, err := ethclient.Dial(cli.NetworkWS)
	if err != nil {
		panic(err)
	}
	httpClient, err := ethclient.Dial(cli.NetworkHTTP)
	return &Client{HttpCli: httpClient, WsCli: wsClient}
}
