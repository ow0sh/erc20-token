package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ow0sh/erc20-token/config"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/usecases"
)

const defaultConfigPath = "./config/config.json"

func main() {
	cfg, err := config.NewConfig(defaultConfigPath)
	if err != nil {
		panic(err)
	}

	log := cfg.Log()
	client := cfg.Client()
	keys := cfg.Keys()

	address := crypto.PubkeyToAddress(*keys.PublicKeyECDSA)

	if keys.ContractAddress.Hex() == "0x0000000000000000000000000000000000000000" {
		addr, err := Deploy(client, keys.PrivateKey, address)
		if err != nil {
			log.Error(err)
		}
		keys.ContractAddress = *addr
		if err = RewriteConfig(keys); err != nil {
			log.Error(err)
		}
		log.Info("Contract successfully deployed")
	}

	instance, err := contracts.NewContracts(keys.ContractAddress, client)
	if err != nil {
		log.Error(err)
	}

	contractUseCase := usecases.NewContractUseCase(instance, keys, client)
	balance, err := contractUseCase.BalanceOf(address)
	if err != nil {
		log.Error(err)
	}
	log.Info(balance)

	toAddr := common.HexToAddress("0x159B2dCdcd6DE5EC249Ca5ed6B5F8dD05B24DA39")
	tx, err := contractUseCase.Transfer(toAddr, big.NewInt(10000))
	if err != nil {
		log.Error(err)
	}

	log.Info(tx.Hash())

	balance, err = contractUseCase.BalanceOf(toAddr)
	if err != nil {
		log.Error(err)
	}
	log.Info(balance)
}
