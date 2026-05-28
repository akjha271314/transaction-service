package repository

import "database/sql"

type TxRunner interface {
	RunInTx(fn func(tx *sql.Tx) error) error
}

type txRunner struct {
	db *sql.DB
}

func NewTxRunner(db *sql.DB) TxRunner {
	return &txRunner{db: db}
}

func (r *txRunner) RunInTx(fn func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
