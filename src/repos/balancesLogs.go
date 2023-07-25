package repos

import "context"

type BalanceLogsTransactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]BalanceLog, error)
	QUpdate(ctx context.Context, q DB) ([]BalanceLog, error)
	QGet(ctx context.Context, q DB) (*BalanceLog, error)
}

type BalanceLogsInserter interface {
	Create(ctx context.Context) ([]BalanceLog, error)
	SetCreate(...CreateBalanceLog) BalanceLogsInserter
}

type BalanceLogsSelector interface {
	FilterByBalanceLogsId(...int64) BalanceLogsSelector
	OrderBy(string, string) BalanceLogsSelector
	WhereAddrS(string) BalanceLogsSelector
	Limit(uint64) BalanceLogsSelector

	Select(ctx context.Context) ([]BalanceLog, error)
	Get(ctx context.Context) (*BalanceLog, error)
}

type BalanceLogsUpdater interface {
	WhereId(...int64) BalanceLogsUpdater
	WhereAddrU(string) BalanceLogsUpdater
	Update(ctx context.Context, column string, value interface{}) ([]BalanceLog, error)
}

type BalanceLogsRepo interface {
	Updater() BalanceLogsUpdater
	Selector() BalanceLogsSelector
	Inserter() BalanceLogsInserter
	Tx() BalanceLogsTransactor
}

type CreateBalanceLog struct {
	Address string
	Balance string
}

type BalanceLog struct {
	Id      int64  `db:"id"`
	Address string `db:"address"`
	Balance string `db:"balance"`
}
