package crud

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/student"
	"github.com/jackc/pgx/v5"
)

type repo struct {
	db pgx.Conn
}

func NewStudentRepo(conn pgx.Conn) student.StudentRepository {
	return &repo{
		db: conn,
	}
}

func (this repo) GetStudents(ctx context.Context, filters student.GetStudentsFilters) ([]student.Student, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	sql, args, err := sq.
		Select(
			"s.uuid",
			"s.first_name",
			"s.last_name",
			"s.middle_name",
			"s.group_id",
			"g.name",
		).
		From("students s").
		LeftJoin("groups g ON s.group_id = g.id").
		Where(sq.Eq{"student_uuid": filters.IDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := this.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var result []student.Student
	for rows.Next() {
		stdnt := student.Student{Group: &student.Group{}}
		if err := rows.Scan(
			&stdnt.UUID,
			&stdnt.FirstName,
			&stdnt.LastName,
			&stdnt.MiddleName,
			&stdnt.Group.ID,
			&stdnt.Group.Name,
		); err != nil {
			return nil, err
		}

		result = append(result, stdnt)
	}

	return result, nil
}
