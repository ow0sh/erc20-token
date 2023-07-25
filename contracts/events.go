package contracts

import "github.com/ethereum/go-ethereum/common"

func (c ContractsTransfer) Topic() common.Hash {
	return common.HexToHash("0x9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc44")
}

func (c ContractsBalanceChanged) Topic() common.Hash {
	return common.HexToHash("0xa448afda7ea1e3a7a10fcab0c29fe9a9dd85791503bf0171f281521551c7ec05")
}
