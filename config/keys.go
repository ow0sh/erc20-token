package config

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type keys struct {
	PrivateKey      string `json:"privateKey"`
	ContractAddress string `json:"contractAddress"`
}

func (k *keys) Keys() Keys {
	contractAddr := common.HexToAddress(k.ContractAddress)
	privateKeyECDSA, err := crypto.HexToECDSA(k.PrivateKey)
	if err != nil {
		panic(err)
	}
	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("failed to get ecdsa public key")
	}
	return Keys{PrivateKey: privateKeyECDSA, PublicKey: publicKey, PublicKeyECDSA: publicKeyECDSA, ContractAddress: contractAddr}
}
