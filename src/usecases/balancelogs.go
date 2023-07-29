package usecases

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ow0sh/erc20-token/src/repos"
	"github.com/pkg/errors"
)

type BalanceLogUseCase struct {
	repo repos.BalanceLogsRepo
}

func NewBalanceLogUseCase(repo repos.BalanceLogsRepo) BalanceLogUseCase {
	return BalanceLogUseCase{repo: repo}
}

func (use BalanceLogUseCase) GetBalance(ctx context.Context, address string) (*repos.BalanceLog, error) {
	balance, err := use.repo.Selector().WhereAddrS(address).Get(ctx)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "failed to get balance")
	}
	return balance, nil
}

func (use BalanceLogUseCase) Exist(ctx context.Context, address string) (bool, error) {
	_, err := use.repo.Selector().WhereAddrS(address).Get(ctx)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "failed to get exist")
	}
	return true, nil
}

func (use BalanceLogUseCase) Insert(ctx context.Context, balance repos.CreateBalanceLog) ([]repos.BalanceLog, error) {
	exist, err := use.Exist(ctx, balance.Address)
	if exist == true {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to get exist")
	}

	balanceDB, err := use.repo.Inserter().SetCreate(balance).Create(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert balance log")
	}

	return balanceDB, nil
}

func (use BalanceLogUseCase) UpdateBalance(ctx context.Context, address string, newBalStr string, operand string) ([]repos.BalanceLog, error) {
	exist, err := use.Exist(ctx, address)
	if exist == false {
		return nil, nil
	}

	oldBalStr, err := use.GetBalance(ctx, address)
	if err != nil {
		return nil, err
	}
	oldbal, _ := strconv.Atoi(*&oldBalStr.Balance)
	newbal, _ := strconv.Atoi(*&newBalStr)
	var balance string
	if operand == "incr" {
		balance = fmt.Sprint(oldbal + newbal)
	}
	if operand == "decr" {
		balance = fmt.Sprint(oldbal - newbal)
	}

	balanceDB, err := use.repo.Updater().WhereAddrU(address).Update(ctx, "balance", balance)
	if err != nil {
		return nil, err
	}

	return balanceDB, err
}
