package postgres

import (
	"context"

	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	db *pgxpool.Pool
}

func NewTransaction(conn *pgxpool.Pool) *Transaction {
	return &Transaction{
		db: conn,
	}
}

func (t *Transaction) PerformTransaction(ctx context.Context, wrapper func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	defer tx.Rollback(ctx)

	contextWithTransaction := context.WithValue(ctx, valueobjects.TransactionKey{}, tx)

	err = wrapper(contextWithTransaction)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
