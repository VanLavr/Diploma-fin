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
	GetExamByID(context.Context, query.GetExamsFilters) (*entities.Exam, error)
	UpdateDebt(context.Context, commands.UpdateDebtByID) error
	CreateExam(context.Context, commands.CreateExam) (int64, error)
	CreateDebt(context.Context, commands.CreateDebt) (int64, error)
	UpdateExam(context.Context, commands.UpdateExamByID) error
	DeleteExam(context.Context, commands.DeleteExam) error
	DeleteDebt(context.Context, commands.DeleteDebt) error
	SearchExams(context.Context, query.SearchExamFilters) ([]entities.Exam, error)
	SearchDebts(context.Context, query.SearchDebtsFilters) ([]entities.Debt, error)
}
