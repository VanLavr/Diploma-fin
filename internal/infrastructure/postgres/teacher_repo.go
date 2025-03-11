package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/entities"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/jackc/pgx/v5"
)

type teacherRepo struct {
	db pgx.Conn
}

func NewTeacherRepo(conn pgx.Conn) repositories.TeacherRepository {
	return &teacherRepo{
		db: conn,
	}
}

func (this teacherRepo) GetTeachers(ctx context.Context, filters query.GetTeachersFilters) ([]entities.Teacher, error) {
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

	var result []entities.Teacher
	for rows.Next() {
		var teacher entities.Teacher
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
