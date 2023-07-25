package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/pkg/ethbackend"
)

func (ch chain) Contracts() *contracts.Contracts {
	return ch.contracts
}

func (ch chain) ChainId() int {
	return ch.chainId
}

func (ch chain) BlockAtStart() uint64 {
	return ch.blockAtStart
}

func (ch chain) ContractAddr() common.Address {
	return ch.contrAddr
}

func (ch chain) ETHCLI() *ethbackend.EthBackend {
	return ch.ethCli
}
