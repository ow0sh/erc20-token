package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ow0sh/erc20-token/src/models"
)

var (
	ErrSuchObjectAlreadyExist = errors.New("such object is already exist")
	ErrNothingUpdated         = errors.New("nothing updated")
)

type DB interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Transactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]int64, error)
	QUpdate(ctx context.Context, q DB) error
}

type OrderBy string

func BalanceLogToCreateRepo(balanceLogs ...models.BalanceLog) []CreateBalanceLog {
	result := make([]CreateBalanceLog, len(balanceLogs))

	for i, balanceLog := range balanceLogs {
		result[i] = CreateBalanceLog{
			Address: balanceLog.Address,
			Balance: balanceLog.Balance,
		}
	}

	return result
}

func TransferLogToCreateRepo(transferLogs ...models.TransferLog) []CreateTransferLog {
	result := make([]CreateTransferLog, len(transferLogs))

	for i, transferLog := range transferLogs {
		result[i] = CreateTransferLog{
			Hash:     transferLog.Hash,
			FromAddr: transferLog.FromAddr,
			ToAddr:   transferLog.ToAddr,
			Value:    fmt.Sprint(transferLog.Value),
			Tokens:   fmt.Sprint(transferLog.Tokens),
		}
	}

	return result
}
