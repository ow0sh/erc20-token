package config

import (
	crypto1 "crypto"
	"crypto/ecdsa"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type Keys struct {
	PrivateKey      *ecdsa.PrivateKey
	PublicKey       crypto1.PublicKey
	PublicKeyECDSA  *ecdsa.PublicKey
	ContractAddress common.Address
}

type Config interface {
	Log() *logrus.Logger
	Client() *ethclient.Client
	Keys() Keys
}

type config struct {
	logger
	client
	keys
}

func NewConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg := config{}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
