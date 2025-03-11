package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	entities "github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type ExamRepository interface {
	GetDebts(context.Context, query.GetDebtsFilters) ([]entities.Debt, error)
	GetExams(context.Context, query.GetExamsFilters) ([]entities.Exam, error)
	UpdateDebt(context.Context, commands.UpdateDebtByID) error
}
