package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

// Transactional represents a transactional object, implementing the dbports.Transactional interface.
type Transactional struct {
	db *sql.DB
}

// NewTransactional creates a new Transactional.
func NewTransactional(db *sql.DB) Transactional {
	return Transactional{
		db: db,
	}
}

// ctxKey is a key for a transaction in a context.
type ctxKey int

const (
	// txKey is the key for a transaction in a context.
	txKey ctxKey = iota
)

// tx creates a new transaction.
func (r Transactional) tx(ctx context.Context) (tx *sql.Tx, err error) {
	if ctx == nil {
		return nil, NilCtxError
	}

	return r.db.BeginTx(ctx, nil)
}

// ctxGetTx returns a transaction from a context if it has been set.
func (r Transactional) ctxGetTx(ctx context.Context) (tx *sql.Tx, err error) {
	if ctx == nil {
		return nil, NilCtxError
	}

	v := ctx.Value(txKey)
	if v == nil {
		return nil, NoTxInCtxError
	}

	tx, ok := v.(*sql.Tx)
	if !ok {
		return nil, InvalidTxInCtxError
	}

	return tx, nil
}

// CtxTx returns a context with a transaction.
func (r Transactional) CtxTx(ctx context.Context) (context.Context, error) {
	if ctx == nil {
		return nil, NilCtxError
	}
	tx, err := r.tx(ctx)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, txKey, tx), nil
}

// TxRollback rolls back a transaction.
func (r Transactional) TxRollback(ctx context.Context) error {
	if ctx == nil {
		return NilCtxError
	}

	tx, err := r.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return nil
		}
		return err
	}

	return tx.Rollback()
}

// TxCommit commits a transaction.
func (r Transactional) TxCommit(ctx context.Context) error {
	if ctx == nil {
		return NilCtxError
	}

	tx, err := r.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return nil
		}
		return err
	}

	return tx.Commit()
}
