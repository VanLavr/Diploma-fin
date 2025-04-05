package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type studentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(conn *pgxpool.Pool) repositories.StudentRepository {
	return &studentRepo{
		db: conn,
	}
}

func (this studentRepo) GetStudents(ctx context.Context, filters query.GetStudentsFilters) ([]models.Student, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	query := sq.Select(
		"s.uuid",
		"s.first_name",
		"s.last_name",
		"s.middle_name",
		"s.group_id",
		"s.email",
		"g.name",
	).From("students s")

	if len(filters.IDs) != 0 {
		query = query.Where(sq.Eq{"uuid": filters.IDs})
	}
	if len(filters.Emails) != 0 {
		query = query.Where(sq.Eq{"email": filters.Emails})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := this.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var result []models.Student
	for rows.Next() {
		stdnt := models.Student{Group: &models.Group{}}
		if err := rows.Scan(
			&stdnt.UUID,
			&stdnt.FirstName,
			&stdnt.LastName,
			&stdnt.MiddleName,
			&stdnt.Group.ID,
			&stdnt.Email,
			&stdnt.Group.Name,
		); err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return nil, err
		}

		result = append(result, stdnt)
	}

	return result, nil
}
