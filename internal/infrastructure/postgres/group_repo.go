package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
	"github.com/VanLavr/Diploma-fin/utils/tools"
)

type groupRepo struct {
	db *pgxpool.Pool
}

// CreateGroup implements repositories.GroupRepository.
func (g *groupRepo) CreateGroup(ctx context.Context, group commands.CreateGroup) (int64, error) {
	sql, args, err := sq.
		Insert("groups").
		SetMap(sq.Eq{
			"name": group.Name,
		}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	row := g.db.QueryRow(ctx, sql, args...)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

// DeleteGroup implements repositories.GroupRepository.
func (g *groupRepo) DeleteGroup(ctx context.Context, group commands.DeleteGroup) error {
	sql, args, err := sq.Delete("groups").Where(sq.Eq{"id": group.ID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if _, err = g.db.Exec(ctx, sql, args...); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

// GetGroupByID implements repositories.GroupRepository.
func (g *groupRepo) GetGroupByID(ctx context.Context, id int64) (*models.Group, error) {
	sql, args, err := sq.Select("id", "name").From("groups").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	row := g.db.QueryRow(ctx, sql, args...)

	result := new(models.Group)
	if err := row.Scan(&result.ID, &result.Name); err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// GetGroups implements repositories.GroupRepository.
func (g *groupRepo) GetGroups(ctx context.Context, filters query.GetGroupsFilters) ([]models.Group, error) {
	sql, args, err := sq.
		Select("id", "name").
		From("groups").
		Where(sq.Eq{"id": filters.IDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := g.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var result []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(&group.ID, &group.Name); err != nil {
			return nil, err
		}

		result = append(result, group)
	}

	return result, nil
}

// UpdateGroup implements repositories.GroupRepository.
func (g *groupRepo) UpdateGroup(ctx context.Context, group commands.UpdateGroup) error {
	query := sq.Update("groups").
		Set("name", group.Name).
		Where(sq.Eq{"id": group.ID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if tx, ok := tools.GetTransaction(ctx); ok {
		_, err = tx.Exec(ctx, sql, args...)
	} else {
		_, err = g.db.Exec(ctx, sql, args...)
	}

	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

func NewGroupRepo(conn *pgxpool.Pool) repositories.GroupRepository {
	return &groupRepo{
		db: conn,
	}
}
