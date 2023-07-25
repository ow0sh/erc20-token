package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/src/models"
	"github.com/pkg/errors"
)

func Deploy(cli *ethclient.Client, private *ecdsa.PrivateKey, address common.Address) (*common.Address, error) {
	nonce, err := cli.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(private)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice, _ = cli.SuggestGasPrice(context.Background())
	auth.GasLimit = uint64(30000000)
	auth.Value = big.NewInt(0)

	contractAddr, tx, _, err := contracts.DeployContracts(auth, cli)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to deploy contract")
	}
	receipt, err := cli.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get tx receipt")
	}
	if receipt.Status == 0 {
		return nil, errors.Wrap(nil, "Tx receipt is 0")
	}

	return &contractAddr, nil
}

func RewriteConfig(keys models.Keys) error {
	fileToRead, err := os.Open("./src/config/config.json")
	defer fileToRead.Close()
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	cfg := models.ConfigModel{}
	if err = json.NewDecoder(fileToRead).Decode(&cfg); err != nil {
		return errors.Wrap(err, "failed to decode file into variable")
	}
	cfg.ContractAddress = keys.ContractAddress.String()
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal variable")
	}

	fileToWrite, err := os.OpenFile("./src/config/config.json", os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to open filetowrite")
	}
	if _, err = fileToWrite.Write(jsonData); err != nil {
		return errors.Wrap(err, "failed to write into file")
	}
	return nil
}
