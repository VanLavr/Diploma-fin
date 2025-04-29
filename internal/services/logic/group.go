package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type GroupUsecase interface {
	CreateGroup(context.Context, types.Group) (int64, error)
	DeleteGroup(context.Context, int64) error
	UpdateGroup(context.Context, types.Group) error
	GetGroupByID(context.Context, int64) (*types.Group, error)
	GetGroups(context.Context, int64, int64) ([]types.Group, error)
}
