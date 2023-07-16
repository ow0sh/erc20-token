package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ow0sh/erc20-token/config"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/usecases"
)

const defaultConfigPath = "./config/config.json"
const contractAddr = "0xBe45e86581dc16f19C0bCc2F46F1ebBf311A385a"

func main() {
	cfg, err := config.NewConfig(defaultConfigPath)
	if err != nil {
		panic(err)
	}

	log := cfg.Log()
	client := cfg.Client()
	keys := cfg.Keys()

	address := crypto.PubkeyToAddress(*keys.PublicKeyECDSA)

	instance, err := contracts.NewContracts(common.HexToAddress(contractAddr), client)
	if err != nil {
		log.Error(err)
	}

	contractUseCase := usecases.NewContractUseCase(instance)
	balance, err := contractUseCase.BalanceOf(address)
	if err != nil {
		log.Error(err)
	}
	log.Info(balance)
}
