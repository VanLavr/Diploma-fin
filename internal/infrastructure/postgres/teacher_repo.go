package postgres

import (
	"context"
	"fmt"

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

type teacherRepo struct {
	db *pgxpool.Pool
}

// SearchTeachers implements repositories.TeacherRepository.
func (this *teacherRepo) SearchTeachers(ctx context.Context, filters query.SearchTeacherFilters) ([]models.Teacher, error) {
	query := sq.Select(
		"uuid",
		"first_name",
		"last_name",
		"middle_name",
		"email",
	).
		From("teachers")

	conditions := sq.And{}
	if len(filters.UUIDs) > 0 {
		conditions = append(conditions, sq.Eq{"uuid": filters.UUIDs})
	}
	if len(filters.FirstNames) > 0 {
		conditions = append(conditions, sq.Eq{"first_name": filters.FirstNames})
	}
	if len(filters.LastNames) > 0 {
		conditions = append(conditions, sq.Eq{"last_name": filters.LastNames})
	}
	if len(filters.MiddleNames) > 0 {
		conditions = append(conditions, sq.Eq{"middle_name": filters.MiddleNames})
	}
	if len(filters.Emails) > 0 {
		conditions = append(conditions, sq.Eq{"email": filters.Emails})
	}

	if len(conditions) > 0 {
		query = query.Where(conditions)
	}

	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := this.db.Query(ctx, sqlQuery, args...)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(
			&teacher.UUID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.MiddleName,
			&teacher.Email,
		); err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teachers, nil
}

func NewTeacherRepo(conn *pgxpool.Pool) repositories.TeacherRepository {
	return &teacherRepo{
		db: conn,
	}
}

// CreateTeacher implements repositories.TeacherRepository.
func (this *teacherRepo) CreateTeacher(ctx context.Context, teacher commands.CreateTeacher) (string, error) {
	sql, args, err := sq.
		Insert("teachers").
		SetMap(sq.Eq{
			"uuid":        sq.Expr("uuid_generate_v4()"),
			"first_name":  teacher.FirstName,
			"last_name":   teacher.LastName,
			"middle_name": teacher.MiddleName,
			"email":       teacher.Email,
			"password":    teacher.Password,
		}).
		Suffix("RETURNING uuid").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	row := this.db.QueryRow(ctx, sql, args...)

	var id string
	if err := row.Scan(&id); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return id, nil
}

// DeleteTeacher implements repositories.TeacherRepository.
func (this *teacherRepo) DeleteTeacher(ctx context.Context, teacher commands.DeleteTeacher) error {
	sql, args, err := sq.Delete("teachers").Where(sq.Eq{"uuid": teacher.UUID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if _, err = this.db.Exec(ctx, sql, args...); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

// GetTeacherByUUID implements repositories.TeacherRepository.
func (this *teacherRepo) GetTeacherByUUID(ctx context.Context, uuid string) (*models.Teacher, error) {
	query := sq.Select(
		"s.uuid",
		"s.first_name",
		"s.last_name",
		"s.middle_name",
		"s.email",
	).From("teachers s")

	query = query.Where(sq.Eq{"s.uuid": uuid})

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := this.db.QueryRow(ctx, sql, args...)

	var result models.Teacher

	if err := row.Scan(
		&result.UUID,
		&result.FirstName,
		&result.LastName,
		&result.MiddleName,
		&result.Email,
	); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	return &result, nil
}

// UpdateTeacher implements repositories.TeacherRepository.
func (this *teacherRepo) UpdateTeacher(ctx context.Context, teacher commands.UpdateTeacher) error {
	query := sq.Update("teachers").
		Set("first_name", teacher.FirstName).
		Set("last_name", teacher.LastName).
		Set("middle_name", teacher.MiddleName).
		Set("email", teacher.Email).
		Where(sq.Eq{"uuid": teacher.UUID}).
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
		"password",
	).From("teachers")

	if len(filters.UUIDs) != 0 {
		query = query.Where(sq.Eq{"uuid": filters.UUIDs})
	}
	if len(filters.Emails) != 0 {
		query = query.Where(sq.Eq{"email": filters.Emails})
	}
	if filters.Limit != 0 {
		query = query.Limit(uint64(filters.Limit))
	}
	if filters.Offset != 0 {
		query = query.Offset(uint64(filters.Offset))
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
			&teacher.Password,
		); err != nil {
			return nil, err
		}

		result = append(result, teacher)
	}

	return result, nil
}
