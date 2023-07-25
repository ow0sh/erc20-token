package sqlx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/erc20-token/src/repos"
)

type blocksRepo struct {
	baseRepo[repos.Block]
}

func NewBlocksRepo(db *sqlx.DB) repos.BlocksRepo {
	return &blocksRepo{
		newBaseRepo[repos.Block](db, "blocks", "value"),
	}
}

func (s blocksRepo) Inserter() repos.BlocksInserter {
	return s
}

func (s blocksRepo) SetCreate(transfers ...repos.CreateTransferLog) repos.BlocksInserter {
	for _, transfer := range transfers {
		s.q.sqlInsert = s.q.sqlInsert.Values(transfer.Hash, transfer.FromAddr, transfer.ToAddr, transfer.Value, transfer.Tokens)
	}
	return s
}

func (s blocksRepo) Create(ctx context.Context) ([]repos.Block, error) {
	return s.baseRepo.Create(ctx)
}

func (s blocksRepo) Selector() repos.BlocksSelector {
	return s
}

func (s blocksRepo) WhereBlockchainS(blockchain string) repos.BlocksSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"blockchain": blockchain})
	return s
}

func (s blocksRepo) FilterByBlocksId(ids ...int64) repos.BlocksSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"id": ids})
	return s
}

func (s blocksRepo) Limit(u uint64) repos.BlocksSelector {
	s.baseRepo = s.baseRepo.Limit(u)
	return s
}

func (s blocksRepo) OrderBy(by string, order string) repos.BlocksSelector {
	s.q.sqlSelect = s.q.sqlSelect.OrderBy(by, order)
	return s
}

func (s blocksRepo) Tx() repos.BlocksTransactor {
	return s
}

func (s blocksRepo) Updater() repos.BlocksUpdater {
	return s
}

func (s blocksRepo) WhereId(ids ...int64) repos.BlocksUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"id": ids})
	return s
}

func (s blocksRepo) WhereBlockchainU(blockchain string) repos.BlocksUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"blockchain": blockchain})
	return s
}
