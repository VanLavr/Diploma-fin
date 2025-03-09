package crud

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
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

func (this repo) GetExams(ctx context.Context, filters exam.GetExamsFilters) ([]exam.Exam, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	sql, args, err := sq.
		Select("id", "name").
		From("exams").
		Where(sq.Eq{"id": filters.IDs}).
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

func (this repo) GetDebts(ctx context.Context, filters exam.GetDebtsFilters) ([]exam.Debt, error) {
	if err := filters.Validate(); err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "invalid filters")
	}

	query := sq.Select(
		"e.exam_id",
		"e.exam_name",
		"t.uuid",
		"t.first_name",
		"t.last_name",
		"t.middle_name",
		"s.uuid",
		"s.first_name",
		"s.last_name",
		"s.middle_name",
		"d.date",
	)
	query = query.From("debts d")

	if len(filters.StudentUUIDs) > 0 {
		query = query.LeftJoin("students s ON d.student_uuid = s.uuid")
		query = query.Where(sq.Eq{"s.student_uuid": filters.StudentUUIDs})
	}
	if len(filters.TeacherUUIDs) > 0 {
		query = query.LeftJoin("teachers t ON d.teacher_uuid = t.uuid")
		query = query.Where(sq.Eq{"t.teacher_uuid": filters.TeacherUUIDs})
	}
	if len(filters.ExamIDs) > 0 {
		query = query.LeftJoin("exams e ON d.exam_id = e.id")
		query = query.Where(sq.Eq{"e.exam_id": filters.ExamIDs})
	}
	query = query.PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not build sql")
	}

	rows, err := this.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not perform query")
	}

	var result []exam.Debt
	for rows.Next() {
		var debt exam.Debt
		if err := rows.Scan(
			&debt.ExamID,
			&debt.ExamName,
			&debt.TeacherUUID,
			&debt.TeacherFirstName,
			&debt.TeacherLastName,
			&debt.TeacherMiddleName,
			&debt.StudentUUID,
			&debt.StudentFirstName,
			&debt.StudentLastName,
			&debt.StudentMiddleName,
			&debt.Date,
		); err != nil {
			return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not scan result")
		}

		result = append(result, debt)
	}

	return result, nil
}
