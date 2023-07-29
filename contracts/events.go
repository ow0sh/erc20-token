package contracts

import "github.com/ethereum/go-ethereum/common"

func (c ContractsTransfer) Topic() common.Hash {
	return common.HexToHash("0x9ed053bb818ff08b8353cd46f78db1f0799f31c9e4458fdb425c10eccd2efc44")
}
