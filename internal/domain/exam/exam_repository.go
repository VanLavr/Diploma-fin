package exam

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Connector interface {
	ConnectToPostgres(string) (pgx.Conn, error)
}

type ExamRepository interface {
	GetAllDebts(context.Context, GetAllDebtsFilter) ([]Exam, error)
}
