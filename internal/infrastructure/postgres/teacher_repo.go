package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
)

type teacherRepo struct {
	db *pgxpool.Pool
}

func NewTeacherRepo(conn *pgxpool.Pool) repositories.TeacherRepository {
	return &teacherRepo{
		db: conn,
	}
}

func (this teacherRepo) GetTeachers(ctx context.Context, filters query.GetTeachersFilters) ([]models.Teacher, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	query := sq.Select(
		"uuid",
		"first_name",
		"last_name",
		"middle_name",
		"email",
	).From("teachers")

	if len(filters.UUIDs) != 0 {
		query = query.Where(sq.Eq{"uuid": filters.UUIDs})
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

	var result []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(
			&teacher.UUID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.MiddleName,
			&teacher.Email,
		); err != nil {
			return nil, err
		}

		result = append(result, teacher)
	}

	return result, nil
}
