package sqlx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/erc20-token/src/repos"
)

type balanceLogRepo struct {
	baseRepo[repos.BalanceLog]
}

func NewBalanceLogRepo(db *sqlx.DB) repos.BalanceLogsRepo {
	return &balanceLogRepo{
		newBaseRepo[repos.BalanceLog](db, "balances", "address", "balance"),
	}
}

func (s *balanceLogRepo) Inserter() repos.BalanceLogsInserter {
	return s
}

func (s *balanceLogRepo) SetCreate(balances ...repos.CreateBalanceLog) repos.BalanceLogsInserter {
	for _, balance := range balances {
		s.q.sqlInsert = s.q.sqlInsert.Values(balance.Address, balance.Balance)
	}
	return s
}

func (s *balanceLogRepo) Create(ctx context.Context) ([]repos.BalanceLog, error) {
	return s.baseRepo.Create(ctx)
}

func (s *balanceLogRepo) Selector() repos.BalanceLogsSelector {
	return s
}

func (s *balanceLogRepo) FilterByBalanceLogsId(ids ...int64) repos.BalanceLogsSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"id": ids})
	return s
}

func (s *balanceLogRepo) Limit(u uint64) repos.BalanceLogsSelector {
	s.baseRepo = s.baseRepo.Limit(u)
	return s
}

func (s *balanceLogRepo) OrderBy(by string, order string) repos.BalanceLogsSelector {
	s.q.sqlSelect = s.q.sqlSelect.OrderBy(by, order)
	return s
}

func (s *balanceLogRepo) WhereAddrS(addr string) repos.BalanceLogsSelector {
	s.q.sqlSelect = s.q.sqlSelect.Where(squirrel.Eq{"address": addr})
	return s
}

func (s *balanceLogRepo) Tx() repos.BalanceLogsTransactor {
	return s
}

func (s *balanceLogRepo) Updater() repos.BalanceLogsUpdater {
	return s
}

func (s *balanceLogRepo) WhereId(ids ...int64) repos.BalanceLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"id": ids})
	return s
}

func (s *balanceLogRepo) WhereAddrU(addr string) repos.BalanceLogsUpdater {
	s.q.sqlUpdate = s.q.sqlUpdate.Where(squirrel.Eq{"address": addr})
	return s
}
