package chain

import (
	"context"

	"github.com/ow0sh/erc20-token/contracts"
	"github.com/pkg/errors"
)

func (ch *chain) init() error {
	var err error

	chainIdBig, err := ch.ethCli.ChainID(context.Background())
	if err != nil {
		return errors.Wrap(err, "Failed to get bigchainid")
	}

	ch.chainId = int(chainIdBig.Int64())

	ch.contracts, err = contracts.NewContracts(ch.contrAddr, ch.ethCli)
	if err != nil {
		return errors.Wrap(err, "failed to create contracts")
	}

	ch.blockAtStart, err = ch.ethCli.BlockNumber(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get current block number")
	}

	return nil
}
