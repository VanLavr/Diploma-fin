package application

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type groupUsecase struct {
	repo repositories.Repository
}

func NewGroupUsecase(repo repositories.Repository) logic.GroupUsecase {
	return &groupUsecase{
		repo: repo,
	}
}

// CreateGroup implements logic.GroupUsecase.
func (g *groupUsecase) CreateGroup(ctx context.Context, group types.Group) (int64, error) {
	id, err := g.repo.CreateGroup(ctx, commands.CreateGroup{
		Name: group.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

// DeleteGroup implements logic.GroupUsecase.
func (g *groupUsecase) DeleteGroup(ctx context.Context, id int64) error {
	if err := g.repo.DeleteGroup(ctx, commands.DeleteGroup{ID: id}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

// GetGroupByID implements logic.GroupUsecase.
func (g *groupUsecase) GetGroupByID(ctx context.Context, id int64) (*types.Group, error) {
	group, err := g.repo.GetGroupByID(ctx, id)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := types.GroupFromDomain(group)

	return &result, nil
}

// GetGroups implements logic.GroupUsecase.
func (g *groupUsecase) GetGroups(ctx context.Context, limit int64, offset int64) ([]types.Group, error) {
	groups, err := g.repo.GetGroups(ctx, query.GetGroupsFilters{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := make([]types.Group, len(groups))
	for i, group := range groups {
		result[i] = types.GroupFromDomain(&group)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// UpdateGroup implements logic.GroupUsecase.
func (g *groupUsecase) UpdateGroup(ctx context.Context, group types.Group) error {
	err := g.repo.UpdateGroup(ctx, commands.UpdateGroup{ID: group.ID, Name: group.Name})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}
