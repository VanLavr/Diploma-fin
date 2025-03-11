package tools

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TransactionKey struct{}

func GetTransaction(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(TransactionKey{}).(pgx.Tx)
	return tx, ok
}
