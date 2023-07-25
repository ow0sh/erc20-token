package usecases

import (
	"context"

	"github.com/ow0sh/erc20-token/src/repos"
	"github.com/pkg/errors"
)

type BalanceLogUseCase struct {
	repo repos.BalanceLogsRepo
}

func NewBalanceLogUseCase(repo repos.BalanceLogsRepo) BalanceLogUseCase {
	return BalanceLogUseCase{repo: repo}
}

func (use BalanceLogUseCase) GetLog(ctx context.Context, address string) (*repos.BalanceLog, error) {
	balanceLog, err := use.repo.Selector().WhereAddrS(address).Get(ctx)
	if err != nil {
		return nil, err
	}

	return balanceLog, nil
}

func (use BalanceLogUseCase) CreateLog(ctx context.Context, balances ...repos.CreateBalanceLog) ([]repos.BalanceLog, error) {
	balanceDB, err := use.repo.Inserter().SetCreate(balances...).Create(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert balance log")
	}

	return balanceDB, nil
}

func (use BalanceLogUseCase) UpdateLog(ctx context.Context, address string, balance string) ([]repos.BalanceLog, error) {
	balanceDB, err := use.repo.Updater().WhereAddrU(address).Update(ctx, "balance", balance)
	if err != nil {
		return nil, err
	}

	return balanceDB, err
}
