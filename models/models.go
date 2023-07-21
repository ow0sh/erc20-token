package models

type ConfigModel struct {
	PrivateKey      string `json:"privateKey"`
	Network         string `json:"network"`
	ContractAddress string `json:"contractAddress"`
}
