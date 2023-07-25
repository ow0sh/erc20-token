package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/pkg/ethbackend"
	"github.com/pkg/errors"
)

type Chain interface {
	Contracts() *contracts.Contracts
	ChainId() int
	BlockAtStart() uint64
	ContractAddr() common.Address
	ETHCLI() *ethbackend.EthBackend
}

type chain struct {
	contracts    *contracts.Contracts
	ethCli       *ethbackend.EthBackend
	chainId      int
	blockAtStart uint64
	contrAddr    common.Address
}

func NewChain(httpCli *ethclient.Client, whCli *ethclient.Client, contrAddr common.Address) Chain {
	ch := &chain{
		ethCli:    ethbackend.NewEthBackend(httpCli, whCli),
		contrAddr: contrAddr,
	}

	if err := ch.init(); err != nil {
		panic(errors.Wrap(err, "failed to init chain"))
	}

	return ch
}
