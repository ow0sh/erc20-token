package repos

import (
	"context"
)

type BlocksTransactor interface {
	QDelete(ctx context.Context, q DB) error
	QCreate(ctx context.Context, q DB) ([]Block, error)
	QUpdate(ctx context.Context, q DB) ([]Block, error)
	QGet(ctx context.Context, q DB) (*Block, error)
}

type BlocksInserter interface {
	Create(ctx context.Context) ([]Block, error)
	SetCreate(...CreateTransferLog) BlocksInserter
}

type BlocksSelector interface {
	FilterByBlocksId(...int64) BlocksSelector
	OrderBy(string, string) BlocksSelector
	Limit(uint64) BlocksSelector
	WhereBlockchainS(string) BlocksSelector

	Select(ctx context.Context) ([]Block, error)
	Get(ctx context.Context) (*Block, error)
}

type BlocksUpdater interface {
	WhereId(...int64) BlocksUpdater
	WhereBlockchainU(string) BlocksUpdater
	Update(ctx context.Context, column string, value interface{}) ([]Block, error)
}

type BlocksRepo interface {
	Updater() BlocksUpdater
	Selector() BlocksSelector
	Inserter() BlocksInserter
	Tx() BlocksTransactor
}

type CreateBlock struct {
	Value      int64
	Blockchain string
}

type Block struct {
	Id         int64  `db:"id"`
	Value      int64  `db:"value"`
	Blockchain string `db:"blockchain"`
}
