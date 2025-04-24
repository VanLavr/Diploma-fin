package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type ExamUsecase interface {
	DeleteExam(context.Context, int64) error
	UpdateExam(context.Context, types.Debt) error
	CreateExam(context.Context, types.Debt) (int64, error)
	GetExams(context.Context, int64, int64) ([]types.Debt, error)
}
