package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type StudentUsecase interface {
	GetAllDebts(context.Context, string) ([]types.Debt, error)
	SendNotification(context.Context, string, int64) error
}
