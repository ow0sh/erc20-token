package sqlx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/erc20-token/src/repos"
)

type transferLogsRepo struct {
	baseRepo[repos.TransferLog]
}

func NewTransferLogsRepo(db *sqlx.DB) repos.TransferLogsRepo {
	return &transferLogsRepo{
		newBaseRepo[repos.TransferLog](db, "transfers", "hash", "fromaddr", "toaddr", "value", "tokens"),
	}
}

func (s transferLogsRepo) Inserter() repos.TransferLogsInserter {
	return s
}

func (s transferLogsRepo) SetCreate(transfers ...repos.CreateTransferLog) repos.TransferLogsInserter {
	for _, transfer := range transfers {
		s.q.sqlInsert = s.q.sqlInsert.Values(transfer.Hash, transfer.FromAddr, transfer.ToAddr, transfer.Value, transfer.Tokens)
	}
	return s
}

func (s transferLogsRepo) Create(ctx context.Context) ([]repos.TransferLog, error) {
	return s.baseRepo.Create(ctx)
}

func (s transferLogsRepo) Selector() repos.TransferLogsSelector {
	return s
}

func (s transferLogsRepo) FilterByTransferLogsId(ids ...int64) repos.TransferLogsSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"id": ids})
	return s
}

func (s transferLogsRepo) Limit(u uint64) repos.TransferLogsSelector {
	s.baseRepo = s.baseRepo.Limit(u)
	return s
}

func (s transferLogsRepo) OrderBy(by string, order string) repos.TransferLogsSelector {
	s.q.sqlSelect = s.q.sqlSelect.OrderBy(by, order)
	return s
}

func (s transferLogsRepo) Tx() repos.TransferLogsTransactor {
	return s
}

func (s transferLogsRepo) Updater() repos.TransferLogsUpdater {
	return s
}

func (s transferLogsRepo) WhereHash(hash string) repos.TransferLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"hash": hash})
	return s
}

func (s transferLogsRepo) WhereFromAddr(fromAddr string) repos.TransferLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"fromaddr": fromAddr})
	return s
}

func (s transferLogsRepo) WhereToAddr(toAddrs string) repos.TransferLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"toaddr": toAddrs})
	return s
}

func (s transferLogsRepo) WhereId(ids ...int64) repos.TransferLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"id": ids})
	return s
}
