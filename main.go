package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ethereum/go-ethereum/crypto"

	_ "github.com/jackc/pgx"
	"github.com/ow0sh/erc20-token/contracts"
	"github.com/ow0sh/erc20-token/src/chain"
	"github.com/ow0sh/erc20-token/src/config"
	"github.com/ow0sh/erc20-token/src/repos/sqlx"
	"github.com/ow0sh/erc20-token/src/state_syncer"
	"github.com/ow0sh/erc20-token/src/usecases"
)

const defaultConfigPath = "./src/config/config.json"

func main() {
	cfg, err := config.NewConfig(defaultConfigPath)
	if err != nil {
		panic(err)
	}

	log := cfg.Log()
	client := cfg.Client()
	keys := cfg.Keys()
	db := cfg.DB()

	ctx, cancel := ctxWithSig()
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			cancel()
		}
	}()

	address := crypto.PubkeyToAddress(*keys.PublicKeyECDSA)

	if keys.ContractAddress.Hex() == "0x0000000000000000000000000000000000000000" {
		addr, err := Deploy(client.WsCli, keys.PrivateKey, address)
		if err != nil {
			log.Error(err)
		}
		keys.ContractAddress = *addr
		if err = RewriteConfig(keys); err != nil {
			log.Error(err)
		}
		log.Info("Contract successfully deployed")
	}

	instance, err := contracts.NewContracts(keys.ContractAddress, client.WsCli)
	if err != nil {
		log.Error(err)
	}

	ethChain := chain.NewChain(
		cfg.Client().HttpCli, cfg.Client().WsCli, keys.ContractAddress,
	)

	transferUseCase := usecases.NewTransferLogUseCase(sqlx.NewTransferLogsRepo(db))
	balanceUseCase := usecases.NewBalanceLogUseCase(sqlx.NewBalanceLogRepo(db))
	contractUseCase := usecases.NewContractUseCase(instance, keys, client.WsCli)
	eventUseCase := usecases.NewEventUseCase(transferUseCase, balanceUseCase, ethChain, *contractUseCase)
	blockUseCase := usecases.NewBlockUseCase(sqlx.NewBlocksRepo(db))

	group := &sync.WaitGroup{}
	state_syncer.NewStateSync(log, ethChain, eventUseCase, blockUseCase, ethChain.BlockAtStart()).Run(ctx, group)
	group.Wait()
}

func ctxWithSig() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	go func() {
		select {
		case <-ch:
			cancel()
		}
	}()

	return ctx, cancel
}
