package usecases

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/src/chain"
	"github.com/ow0sh/erc20-token/src/repos"
	"github.com/pkg/errors"
)

type EventUseCase interface {
	Process(ctx context.Context, log types.Log) error
}

func NewEventUseCase(transferUse TransferLogUseCase, balanceUse BalanceLogUseCase, ch chain.Chain, contractUse ContractUseCase) EventUseCase {
	return eventUseCase{transferUse: transferUse, balanceUse: balanceUse, chain: ch, contractUse: contractUse}
}

type eventUseCase struct {
	chain       chain.Chain
	transferUse TransferLogUseCase
	balanceUse  BalanceLogUseCase
	contractUse ContractUseCase
}

func (e eventUseCase) Process(ctx context.Context, log types.Log) error {
	switch log.Topics[0] {
	case contracts.ContractsTransfer{}.Topic():
		event, err := e.chain.Contracts().ParseTransfer(log)
		if err != nil {
			return errors.Wrap(err, "failed to parse event")
		}
		createLog := repos.CreateTransferLog{Hash: event.Topic().Hex(), FromAddr: event.From.Hex(),
			ToAddr: event.To.Hex(), Value: fmt.Sprint(event.Value), Tokens: fmt.Sprint(event.Tokens)}
		_, err = e.transferUse.CreateLog(ctx, createLog)
		if err != nil {
			return errors.Wrap(err, "failed to create log")
		}
		break
	case contracts.ContractsBalanceChanged{}.Topic():
		event, err := e.chain.Contracts().ParseBalanceChanged(log)
		if err != nil {
			return errors.Wrap(err, "failed to parse balance event")
		}
		_, err = e.balanceUse.UpdateLog(ctx, event.Addr.Hex(), event.Balance.String())
		if err != nil {
			return errors.Wrap(err, "failed to update balance")
		}
		break
	default:
		return errors.New("such log topic is not exists")
	}
	return nil
}
