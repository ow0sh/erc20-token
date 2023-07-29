package listener

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/pkg/ethbackend"
	"github.com/ow0sh/erc20-token/src/chain"
	"github.com/ow0sh/erc20-token/src/usecases"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type EventListener struct {
	log      *logrus.Entry
	ethcli   *ethbackend.EthBackend
	blockUse usecases.BlockUseCase
	eventUse usecases.EventUseCase
	chain    chain.Chain

	startAt uint64
	net     int

	blockLock    sync.Mutex
	currentBlock uint64
}

var topics = [][]common.Hash{
	{
		contracts.ContractsTransfer{}.Topic(),
	},
}

func NewEventListener(log *logrus.Logger,
	ch chain.Chain,
	eventUse usecases.EventUseCase,
	blockUse usecases.BlockUseCase,
	startAt uint64) *EventListener {
	return &EventListener{log: log.WithFields(logrus.Fields{"listener": "event",
		"network": ch.ChainId()}),
		ethcli:   ch.ETHCLI(),
		eventUse: eventUse,
		blockUse: blockUse,
		chain:    ch,
		startAt:  startAt,
		net:      ch.ChainId()}
}

func (d *EventListener) Start(ctx context.Context, group *sync.WaitGroup) {
	startAt, _ := d.blockUse.GetLastBlock(ctx, fmt.Sprint(d.net))
	if startAt != nil {
		d.startAt = uint64(startAt.Value)
	}
	var err error
	d.currentBlock, err = d.ethcli.BlockNumber(ctx)
	if err != nil {
		panic(err)
	}

	if d.currentBlock < d.startAt {
		d.startAt = d.currentBlock
		d.blockUse.SetLastBlock(ctx, d.startAt, fmt.Sprint(d.net))
	}
	group.Add(2)
	go paginateWithRetries(ctx, d.log, time.Second*3, 20, &d.startAt, group, d.pageLogs)
	go watchWithRetries(ctx, d.log, time.Second, 20, group, d.listenHead)
}

func (d *EventListener) pageLogs(ctx context.Context, startAt *uint64) error {
	startAtBig := new(big.Int)
	if startAt != nil {
		startAtBig = startAtBig.SetUint64(*startAt)
	}
	d.blockLock.Lock()
	toBlock := new(big.Int).SetUint64(d.currentBlock)
	d.blockLock.Unlock()

	if startAtBig.Cmp(toBlock) >= 0 {
		return nil
	}

	if toBlock.Uint64()-startAtBig.Uint64() > 1000 {
		toBlock.SetUint64(startAtBig.Uint64() + 1000)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, time.Minute)
	sub, err := d.ethcli.FilterLogs(timeoutCtx, ethereum.FilterQuery{
		FromBlock: startAtBig,
		ToBlock:   toBlock,
		Addresses: []common.Address{d.chain.ContractAddr()},
		Topics:    topics,
	})
	cancel()
	d.log.WithFields(logrus.Fields{"start_at": *startAt, "end_at": toBlock.Uint64()}).Info("Page")
	if err != nil {
		return errors.Wrap(err, "failed to filter")
	}

	for ind := range sub {
		if isDone(ctx) {
			return errContextIsDone
		}

		log := d.log.WithFields(logrus.Fields{
			"block_num": sub[ind].BlockNumber,
			"tx_hash":   sub[ind].TxHash,
		})

		if err = d.eventUse.Process(ctx, sub[ind], log); err != nil {
			log.WithError(err).Error("failed to process")
		}
		setUint64P(startAt, sub[ind].BlockNumber)
	}
	d.blockUse.SetLastBlock(ctx, *startAt, fmt.Sprint(d.net))
	setUint64P(startAt, toBlock.Uint64())
	return nil
}

func (d *EventListener) listenHead(ctx context.Context) error {
	sink := make(chan *types.Header, 10)
	timeoutedCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	sub, err := d.ethcli.SubscribeNewHead(timeoutedCtx, sink)
	cancel()
	if err != nil {
		return errors.Wrap(err, "failed to subscribe to new head")
	}

	for {
		select {
		case <-ctx.Done():
			d.log.Debug("gracefully shutdown")
			return errContextIsDone
		case err = <-sub.Err():
			return errors.Wrap(err, "header subscription error")
		case h := <-sink:
			d.blockLock.Lock()
			d.currentBlock = h.Number.Uint64()
			d.log.WithField("current_block", d.currentBlock).Info("current block updated")
			d.blockLock.Unlock()
		}
	}
}

func setUint64P(to *uint64, from uint64) {
	if to == nil {
		val := reflect.ValueOf(to)
		val.SetUint(from)
		return
	}

	*to = from
}
