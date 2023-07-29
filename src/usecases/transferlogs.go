package usecases

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/ow0sh/erc20-token/src/repos"
)

type TransferLogUseCase struct {
	repo repos.TransferLogsRepo
}

func NewTransferLogUseCase(repo repos.TransferLogsRepo) TransferLogUseCase {
	return TransferLogUseCase{repo: repo}
}

func (use TransferLogUseCase) Insert(ctx context.Context, transfers ...repos.CreateTransferLog) ([]repos.TransferLog, error) {
	transferDB, err := use.repo.Inserter().SetCreate(transfers...).Create(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert transfer logs")
	}
	return transferDB, nil
}

func (use TransferLogUseCase) Exist(ctx context.Context, hash string) (bool, error) {
	_, err := use.repo.Selector().WhereHash(hash).Get(ctx)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return false, errors.Wrap(err, "failed to get log")
	}
	return true, nil
}
