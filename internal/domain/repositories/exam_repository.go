package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/entities"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/jackc/pgx/v5"
)

type Connector interface {
	ConnectToPostgres(string) (pgx.Conn, error)
}

type ExamRepository interface {
	GetDebts(context.Context, query.GetDebtsFilters) ([]entities.Debt, error)
	GetExams(context.Context, query.GetExamsFilters) ([]entities.Exam, error)
	UpdateDebt(context.Context, commands.UpdateDebtByID) error
}
