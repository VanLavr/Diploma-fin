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

type examRepo struct {
	db *pgxpool.Pool
}

func NewExamRepo(conn *pgxpool.Pool) repositories.ExamRepository {
	return &examRepo{
		db: conn,
	}
}

// CreateDebt implements repositories.ExamRepository.
func (this *examRepo) CreateDebt(ctx context.Context, createDebt commands.CreateDebt) (int64, error) {
	sql, args, err := sq.
		Insert("debts").
		SetMap(sq.Eq{
			"exam_id":      createDebt.ExamID,
			"student_uuid": createDebt.StudentUUID,
			"teacher_uuid": createDebt.TeacherUUID,
		}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	row := this.db.QueryRow(ctx, sql, args...)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

// DeleteDebt implements repositories.ExamRepository.
func (this *examRepo) DeleteDebt(ctx context.Context, debt commands.DeleteDebt) error {
	sql, args, err := sq.Delete("debts").Where(sq.Eq{"id": debt.ID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if _, err = this.db.Exec(ctx, sql, args...); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

// GetExamByID implements repositories.ExamRepository.
func (this *examRepo) GetExamByID(ctx context.Context, query query.GetExamsFilters) (*models.Exam, error) {
	sql, args, err := sq.Select("id", "name").From("exams").Where(sq.Eq{"id": query.IDs}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	row := this.db.QueryRow(ctx, sql, args...)

	result := new(models.Exam)
	if err := row.Scan(&result.ID, &result.Name); err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// DeleteExam implements repositories.ExamRepository.
func (this *examRepo) DeleteExam(ctx context.Context, exam commands.DeleteExam) error {
	sql, args, err := sq.Delete("exams").Where(sq.Eq{"id": exam.ID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if _, err = this.db.Exec(ctx, sql, args...); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

// UpdateExam implements repositories.ExamRepository.
func (this *examRepo) UpdateExam(ctx context.Context, exam commands.UpdateExamByID) error {
	query := sq.Update("exams").
		Set("name", exam.Name).
		Where(sq.Eq{"id": exam.ID}).
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

// CreateExam implements repositories.ExamRepository.
func (this *examRepo) CreateExam(ctx context.Context, exam commands.CreateExam) (int64, error) {
	sql, args, err := sq.
		Insert("exams").
		SetMap(sq.Eq{
			"name": exam.Name,
		}).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	row := this.db.QueryRow(ctx, sql, args...)

	var id int64
	if err := row.Scan(&id); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func (this examRepo) GetExams(ctx context.Context, filters query.GetExamsFilters) ([]models.Exam, error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}
	query := sq.Select("id", "name")
	query = query.From("exams")
	query = query.PlaceholderFormat(sq.Dollar)

	if len(filters.IDs) != 0 {
		query = query.Where(sq.Eq{"id": filters.IDs})
	}
	if filters.Offset != 0 {
		query = query.Offset(uint64(filters.Offset))
	}
	if filters.Limit != 0 {
		query = query.Limit(uint64(filters.Limit))
	}

	sql, args, err := query.ToSql()
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
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "invalid filters")
	}

	query := sq.Select(
		"d.id",
		"e.id",
		"e.name",
		"t.uuid",
		"t.first_name",
		"t.last_name",
		"t.middle_name",
		"t.email",
		"s.uuid",
		"s.first_name",
		"s.last_name",
		"s.middle_name",
		"g.id",
		"g.name",
		"s.email",
		"d.date",
	)
	query = query.From("debts d")
	query = query.LeftJoin("students s ON d.student_uuid = s.uuid")
	query = query.LeftJoin("groups g ON s.group_id = g.id")
	query = query.LeftJoin("teachers t ON d.teacher_uuid = t.uuid")
	query = query.LeftJoin("exams e ON d.exam_id = e.id")

	if filters.Limit != 0 {
		query = query.Limit(uint64(filters.Limit))
	}
	if filters.Offset != 0 {
		query = query.Offset(uint64(filters.Offset))
	}
	if len(filters.StudentUUIDs) > 0 {
		query = query.Where(sq.Eq{"s.uuid": filters.StudentUUIDs})
	}
	if len(filters.TeacherUUIDs) > 0 {
		query = query.Where(sq.Eq{"t.uuid": filters.TeacherUUIDs})
	}
	if len(filters.ExamIDs) > 0 {
		query = query.Where(sq.Eq{"e.id": filters.ExamIDs})
	}
	if len(filters.DebtIDs) > 0 {
		query = query.Where(sq.Eq{"d.id": filters.DebtIDs})
	}
	query = query.PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()

	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not build sql")
	}

	rows, err := this.db.Query(ctx, sql, args...)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "can not perform query")
	}

	var result []models.Debt
	for rows.Next() {
		debt := models.Debt{
			Exam: &models.Exam{},
			Student: &models.Student{
				Group: &models.Group{},
			},
			Teacher: &models.Teacher{},
		}
		if err := rows.Scan(
			&debt.ID,
			&debt.Exam.ID,
			&debt.Exam.Name,
			&debt.Teacher.UUID,
			&debt.Teacher.FirstName,
			&debt.Teacher.LastName,
			&debt.Teacher.MiddleName,
			&debt.Teacher.Email,
			&debt.Student.UUID,
			&debt.Student.FirstName,
			&debt.Student.LastName,
			&debt.Student.MiddleName,
			&debt.Student.Group.ID,
			&debt.Student.Group.Name,
			&debt.Student.Email,
			&debt.Date,
		); err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
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
