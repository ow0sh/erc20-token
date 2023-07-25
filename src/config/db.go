package config

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type db struct {
	dbParams `json:"db"`
	db       *sqlx.DB
}

type dbParams struct {
	URL    string `json:"url"`
	Driver string `json:"driver"`
}

func (d *db) DB() *sqlx.DB {
	if d.db == nil {
		db, err := sqlx.Open(d.Driver, d.URL)
		if err != nil {
			panic(err)
		}

		if err = db.Ping(); err != nil {
			panic(err)
		}

		d.db = db
	}

	return d.db
}
