package crud

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
	"github.com/jackc/pgx/v5"
)

type repo struct {
	db pgx.Conn
}

func NewExamRepo(conn pgx.Conn) exam.ExamRepository {
	return &repo{
		db: conn,
	}
}

func (this repo) GetAllDebts(ctx context.Context, filters exam.GetAllDebtsFilter) ([]exam.Exam, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	sql, args, err := sq.
		Select("id", "name").
		From("debts").
		Where(sq.Eq{"student_uuid": filters.StudentUUID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := this.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	var result []exam.Exam
	for rows.Next() {
		var exam exam.Exam
		if err := rows.Scan(&exam.ID, &exam.Name); err != nil {
			return nil, err
		}

		result = append(result, exam)
	}

	return result, nil
}
