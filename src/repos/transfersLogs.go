package repos

import (
	"context"
)

type TransferLogsTransactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]TransferLog, error)
	QUpdate(ctx context.Context, q DB) ([]TransferLog, error)
	QGet(ctx context.Context, q DB) (*TransferLog, error)
}

type TransferLogsInserter interface {
	Create(ctx context.Context) ([]TransferLog, error)
	SetCreate(...CreateTransferLog) TransferLogsInserter
}

type TransferLogsSelector interface {
	FilterByTransferLogsId(...int64) TransferLogsSelector
	OrderBy(string, string) TransferLogsSelector
	Limit(uint64) TransferLogsSelector

	Select(ctx context.Context) ([]TransferLog, error)
	Get(ctx context.Context) (*TransferLog, error)
}

type TransferLogsUpdater interface {
	WhereId(...int64) TransferLogsUpdater
	WhereHash(base string) TransferLogsUpdater
	WhereFromAddr(fromAddr string) TransferLogsUpdater
	WhereToAddr(toAddr string) TransferLogsUpdater
	Update(ctx context.Context, column string, value interface{}) ([]TransferLog, error)
}

type TransferLogsRepo interface {
	Updater() TransferLogsUpdater
	Selector() TransferLogsSelector
	Inserter() TransferLogsInserter
	Tx() TransferLogsTransactor
}

type CreateTransferLog struct {
	Hash     string
	FromAddr string
	ToAddr   string
	Value    string
	Tokens   string
}

type TransferLog struct {
	Id       int64  `db:"id"`
	Hash     string `db:"hash"`
	FromAddr string `db:"fromaddr"`
	ToAddr   string `db:"toaddr"`
	Value    string `db:"value"`
	Tokens   string `db:"tokens"`
}
