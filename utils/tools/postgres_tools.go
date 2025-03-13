package tools

import (
	"context"

	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/jackc/pgx/v5"
)

func GetTransaction(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(valueobjects.TransactionKey{}).(pgx.Tx)
	return tx, ok
}
