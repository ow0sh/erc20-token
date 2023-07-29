package usecases

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/src/chain"
	"github.com/ow0sh/erc20-token/src/repos"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EventUseCase interface {
	Process(ctx context.Context, log types.Log, logrus *logrus.Entry) error
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

func (e eventUseCase) Process(ctx context.Context, log types.Log, logrus *logrus.Entry) error {
	switch log.Topics[0] {
	case contracts.ContractsTransfer{}.Topic():
		exist, err := e.transferUse.Exist(ctx, log.TxHash.Hex())
		if err != nil {
			return errors.Wrap(err, "failed to get exist")
		}
		if exist == true {
			return nil
		}
		logrus.Info("got event")
		event, err := e.chain.Contracts().ParseTransfer(log)
		if err != nil {
			return errors.Wrap(err, "failed to parse event")
		}

		createTransfer := repos.CreateTransferLog{Hash: log.TxHash.Hex(), FromAddr: event.From.Hex(),
			ToAddr: event.To.Hex(), Value: fmt.Sprint(event.Value), Tokens: fmt.Sprint(event.Tokens)}
		_, err = e.transferUse.Insert(ctx, createTransfer)
		if err != nil {
			return errors.Wrap(err, "failed to create log")
		}

		exist, err = e.balanceUse.Exist(context.Background(), event.From.Hex())
		if err != nil {
			return errors.Wrap(err, "failed to get balance exist")
		}
		if exist == true {
			_, err := e.balanceUse.UpdateBalance(ctx, event.From.Hex(), fmt.Sprint(event.Tokens), "decr")
			if err != nil {
				return errors.Wrap(err, "failed to update balance")
			}
		} else {
			bal, err := e.contractUse.BalanceOf(event.From)
			if err != nil {
				return errors.Wrap(err, "failed to get balance of from")
			}
			createBalance := repos.CreateBalanceLog{Address: event.From.Hex(), Balance: fmt.Sprint(bal)}
			_, err = e.balanceUse.Insert(ctx, createBalance)
			if err != nil {
				return errors.Wrap(err, "failed to insert new balance")
			}
		}

		exist, err = e.balanceUse.Exist(context.Background(), event.To.Hex())
		if err != nil {
			return errors.Wrap(err, "failed to get balance exist")
		}
		if exist == true {
			_, err := e.balanceUse.UpdateBalance(ctx, event.To.Hex(), fmt.Sprint(event.Tokens), "incr")
			if err != nil {
				return errors.Wrap(err, "failed to update balance")
			}
		} else {
			bal, err := e.contractUse.BalanceOf(event.To)
			if err != nil {
				return errors.Wrap(err, "failed to get balance of to")
			}
			createBalance := repos.CreateBalanceLog{Address: event.To.Hex(), Balance: fmt.Sprint(bal)}
			_, err = e.balanceUse.Insert(ctx, createBalance)
			if err != nil {
				return errors.Wrap(err, "failed to insert new balance")
			}
		}

		return nil
	default:
		return errors.New("such log topic is not exists")
	}
}
