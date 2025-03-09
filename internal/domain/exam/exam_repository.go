package exam

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Connector interface {
	ConnectToPostgres(string) (pgx.Conn, error)
}

type ExamRepository interface {
	GetDebts(context.Context, GetDebtsFilters) ([]Debt, error)
	GetExams(context.Context, GetExamsFilters) ([]Exam, error)
}
