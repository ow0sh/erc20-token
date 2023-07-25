package usecases

import (
	"context"

	"github.com/ow0sh/erc20-token/src/repos"
)

type BlockUseCase struct {
	repo repos.BlocksRepo
}

func NewBlockUseCase(repo repos.BlocksRepo) BlockUseCase {
	return BlockUseCase{repo: repo}
}

func (s BlockUseCase) GetLastBlock(ctx context.Context, chain string) (*repos.Block, error) {
	block, err := s.repo.Selector().WhereBlockchainS(chain).Get(ctx)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (s BlockUseCase) SetLastBlock(ctx context.Context, block uint64, blockchain string) error {
	_, err := s.repo.Updater().WhereBlockchainU(blockchain).Update(ctx, "value", block)
	if err != nil {
		return err
	}

	return nil
}
