package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
	"github.com/VanLavr/Diploma-fin/pkg/tools"
	"github.com/jackc/pgx/v5/pgxpool"
)

type examRepo struct {
	db *pgxpool.Pool
}

func NewExamRepo(conn *pgxpool.Pool) repositories.ExamRepository {
	return &examRepo{
		db: conn,
	}
}

func (this examRepo) GetExams(ctx context.Context, filters query.GetExamsFilters) ([]models.Exam, error) {
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

	var result []models.Exam
	for rows.Next() {
		var exam models.Exam
		if err := rows.Scan(&exam.ID, &exam.Name); err != nil {
			return nil, err
		}

		result = append(result, exam)
	}

	return result, nil
}

func (this examRepo) GetDebts(ctx context.Context, filters query.GetDebtsFilters) ([]models.Debt, error) {
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
	if len(filters.DebtIDs) > 0 {
		query = query.Where(sq.Eq{"d.id": filters.DebtIDs})
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

	var result []models.Debt
	for rows.Next() {
		var debt models.Debt
		if err := rows.Scan(
			&debt.Exam.ID,
			&debt.Exam.Name,
			&debt.Teacher.UUID,
			&debt.Teacher.FirstName,
			&debt.Teacher.LastName,
			&debt.Teacher.MiddleName,
			&debt.Student.UUID,
			&debt.Student.FirstName,
			&debt.Student.LastName,
			&debt.Student.MiddleName,
			&debt.Date,
		); err != nil {
			return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not scan result")
		}

		result = append(result, debt)
	}

	return result, nil
}

func (this examRepo) UpdateDebt(ctx context.Context, setCommand commands.UpdateDebtByID) error {
	query := sq.Update("debts").
		Set("date", setCommand.Date).
		Set("teacher_uuid", setCommand.TeacherUUID).
		Set("student_uuid", setCommand.StudentUUID).
		Where(sq.Eq{"id": setCommand.DebtID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if tx, ok := tools.GetTransaction(ctx); ok {
		_, err = tx.Exec(ctx, sql, args...)
	} else {
		_, err = this.db.Exec(ctx, sql, args...)
	}

	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
