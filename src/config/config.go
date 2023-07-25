package config

import (
	"encoding/json"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/ow0sh/erc20-token/src/models"
	"github.com/sirupsen/logrus"
)

type Config interface {
	Log() *logrus.Logger
	Client() *Client
	Keys() models.Keys
	DB() *sqlx.DB
}

type config struct {
	logger
	client
	keys
	db
}

func NewConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}

	cfg := config{}
	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
