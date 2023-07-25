package usecases

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ow0sh/erc20-token/src/repos"
)

type TransferLogUseCase struct {
	repo repos.TransferLogsRepo
}

func NewTransferLogUseCase(repo repos.TransferLogsRepo) TransferLogUseCase {
	return TransferLogUseCase{repo: repo}
}

func (use TransferLogUseCase) CreateLog(ctx context.Context, transfers ...repos.CreateTransferLog) ([]repos.TransferLog, error) {
	transferDB, err := use.repo.Inserter().SetCreate(transfers...).Create(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert transfer logs")
	}
	return transferDB, nil
}
