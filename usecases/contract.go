package usecases

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ow0sh/erc20-token/config"
	"github.com/ow0sh/erc20-token/contracts"
)

type ContractUseCase struct {
	instance *contracts.Contracts
	keys     config.Keys
	client   *ethclient.Client
}

func NewContractUseCase(instance *contracts.Contracts, keys config.Keys, client *ethclient.Client) *ContractUseCase {
	return &ContractUseCase{instance: instance, keys: keys, client: client}
}

func (contr *ContractUseCase) TotalSupply() (*big.Int, error) {
	result, err := contr.instance.TotalSupply(nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (contr *ContractUseCase) BalanceOf(addr common.Address) (*big.Int, error) {
	result, err := contr.instance.BalanceOf(nil, addr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (contr *ContractUseCase) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	auth := bind.NewKeyedTransactor(contr.keys.PrivateKey)
	auth.GasLimit = uint64(30000000)
	auth.GasPrice, _ = contr.client.SuggestGasPrice(context.Background())
	auth.Value = big.NewInt(0)
	fromAddr := crypto.PubkeyToAddress(*contr.keys.PublicKeyECDSA)
	nonce, err := contr.client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err := contr.instance.Transfer(auth, to, value)
	if err != nil {
		return nil, err
	}
	return tx, nil
}
