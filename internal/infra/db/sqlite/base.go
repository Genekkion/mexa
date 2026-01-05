package sqlite

import (
	"context"
	"database/sql"
	"errors"
)

// BaseRepo is the base repository for all repositories, based on pgxpool.Pool.
type BaseRepo struct {
	Transactional
	db *sql.DB
}

// NewBaseRepo creates a new BaseRepo.
func NewBaseRepo(db *DB) BaseRepo {
	return BaseRepo{
		db:            db.db,
		Transactional: NewTransactional(db.db),
	}
}

type RowInterface interface {
	Scan(dest ...any) error
}
type ErrRow struct {
	err error
}

// Scan returns the error.
func (e ErrRow) Scan(_ ...any) error {
	return e.err
}

// QueryRow executes a query expected to return at most one row. It uses the
// transaction in the context if available, otherwise it uses the pool directly.
func (repo *BaseRepo) QueryRow(ctx context.Context, stmt string, args ...any) (row RowInterface) {
	tx, err := repo.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return repo.db.QueryRowContext(ctx, stmt, args...)
		}

		return ErrRow{
			err: err,
		}
	}

	return tx.QueryRowContext(ctx, stmt, args...)
}

// Query executes a query that returns rows. It uses the transaction in the
// context if available, otherwise it uses the pool directly.
func (repo *BaseRepo) Query(ctx context.Context, stmt string, args ...any) (rows *sql.Rows, err error) {
	tx, err := repo.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return repo.db.QueryContext(ctx, stmt, args...)
		}
		return nil, err
	}

	return tx.QueryContext(ctx, stmt, args...)
}

// Exec executes an SQL statement. It uses the transaction in the context if
// available, otherwise it uses the pool directly.
func (repo *BaseRepo) Exec(ctx context.Context, stmt string, args ...any) (result sql.Result, err error) {
	tx, err := repo.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return repo.db.ExecContext(ctx, stmt, args...)
		}

		return nil, err
	}
	return tx.ExecContext(ctx, stmt, args...)
}

func (repo *BaseRepo) Prepare(ctx context.Context, stmt string) (statement *sql.Stmt, err error) {
	tx, err := repo.ctxGetTx(ctx)
	if err != nil {
		if errors.Is(err, NoTxInCtxError) {
			return repo.db.PrepareContext(ctx, stmt)
		}

		return nil, err
	}

	return tx.PrepareContext(ctx, stmt)
}
