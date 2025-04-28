package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type ExamUsecase interface {
	DeleteExam(context.Context, int64) error
	UpdateExam(context.Context, types.Exam) error
	CreateExam(context.Context, types.Exam) (int64, error)
	GetExams(context.Context, int64, int64) ([]types.Exam, error)
	GetExam(context.Context, int64) (*types.Exam, error)

	DeleteDebt(context.Context, int64) error
	UpdateDebt(context.Context, types.Debt) error
	CreateDebt(context.Context, types.Debt) (int64, error)
	GetDebt(context.Context, int64) (*types.Debt, error)
	GetDebts(context.Context, int64, int64) ([]types.Debt, error)
}
