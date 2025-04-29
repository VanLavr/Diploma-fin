package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type GroupRepository interface {
	GetGroupByID(context.Context, int64) (*models.Group, error)
	GetGroups(context.Context, query.GetGroupsFilters) ([]models.Group, error)
	CreateGroup(context.Context, commands.CreateGroup) (int64, error)
	UpdateGroup(context.Context, commands.UpdateGroup) error
	DeleteGroup(context.Context, commands.DeleteGroup) error
}
