package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
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

	sql, args, err := sq.
		Select(
			"uuid",
			"first_name",
			"last_name",
			"middle_name",
			"email",
		).
		From("teachers").
		Where(sq.Eq{"id": filters.UUIDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
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
