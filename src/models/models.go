package models

import (
	"crypto"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ConfigModel struct {
	PrivateKey      string   `json:"privateKey"`
	NetworkWS       string   `json:"networkWS"`
	NetworkHTTP     string   `json:"networkHTTP"`
	ContractAddress string   `json:"contractAddress"`
	DB              dbParams `json:"db"`
}

type dbParams struct {
	URL    string `json:"url"`
	Driver string `json:"driver"`
}

type Keys struct {
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       crypto.PublicKey
	PublicKeyECDSA  *ecdsa.PublicKey
	ContractAddress common.Address
}

type BalanceLog struct {
	Address string
	Balance string
}

type TransferLog struct {
	Hash     string
	FromAddr string
	ToAddr   string
	Value    *big.Int
	Tokens   *big.Int
}
