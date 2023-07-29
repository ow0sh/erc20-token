package state_syncer

import (
	"context"
	"sync"

	"github.com/ow0sh/erc20-token/src/chain"
	"github.com/ow0sh/erc20-token/src/state_syncer/listener"
	"github.com/ow0sh/erc20-token/src/usecases"
	"github.com/sirupsen/logrus"
)

type StateSync struct {
	eventListener *listener.EventListener
	group         *sync.WaitGroup
}

func NewStateSync(log *logrus.Logger,
	ch chain.Chain,
	eventUse usecases.EventUseCase,
	blockUse usecases.BlockUseCase,
	startAt uint64) StateSync {
	return StateSync{eventListener: listener.NewEventListener(log, ch, eventUse, blockUse, startAt)}
}

func (s StateSync) Run(ctx context.Context, group *sync.WaitGroup) {
	s.eventListener.Start(ctx, group)
}
