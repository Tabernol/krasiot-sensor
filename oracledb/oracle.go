package oracledb

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
)

func InitOracle(cfg OracleConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		`user=%s password=%s connectString="%s?wallet_location=%s" libDir=%s`,
		cfg.User, cfg.Password, cfg.ConnStr, cfg.Wallet, cfg.LibDir)

	db, err := sql.Open("godror", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
