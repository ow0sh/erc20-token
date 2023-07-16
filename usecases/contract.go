package usecases

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ow0sh/erc20-token/contracts"
)

type ContractUseCase struct {
	instance *contracts.Contracts
}

func NewContractUseCase(instance *contracts.Contracts) *ContractUseCase {
	return &ContractUseCase{instance: instance}
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
	tx, err := contr.instance.Transfer(nil, to, value)
	if err != nil {
		return nil, err
	}
	return tx, err
}
