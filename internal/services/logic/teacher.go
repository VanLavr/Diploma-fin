package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type TeacherUsecase interface {
	SetDate(context.Context, string, string, int64) error
	GetAllDebts(context.Context, string) ([]types.Debt, error)
}
