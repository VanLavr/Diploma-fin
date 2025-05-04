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

type studentRepo struct {
	db *pgxpool.Pool
}

// ChangeStudentPassword implements repositories.StudentRepository.
func (this *studentRepo) ChangeStudentPassword(ctx context.Context, uuid, password string) error {
	query := sq.Update("students").
		Set("password", password).
		Where(sq.Eq{"uuid": uuid}).
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

// SearchStudents implements repositories.StudentRepository.
func (this *studentRepo) SearchStudents(ctx context.Context, filters query.SearchStudentFilters) ([]models.Student, error) {
	query := sq.Select(
		"uuid",
		"first_name",
		"last_name",
		"middle_name",
		"email",
		"group_id",
	).
		From("students")

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

	var students []models.Student
	for rows.Next() {
		var student models.Student
		var groupID int64
		if err := rows.Scan(
			&student.UUID,
			&student.FirstName,
			&student.LastName,
			&student.MiddleName,
			&student.Email,
			&groupID,
		); err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		// You might want to fetch the group details here if needed
		student.Group = &models.Group{ID: groupID}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return students, nil
}

func NewStudentRepo(conn *pgxpool.Pool) repositories.StudentRepository {
	return &studentRepo{
		db: conn,
	}
}

// GetStudentByUUID implements repositories.StudentRepository.
func (this *studentRepo) GetStudentByUUID(ctx context.Context, uuid string) (*models.Student, error) {
	query := sq.Select(
		"s.uuid",
		"s.first_name",
		"s.last_name",
		"s.middle_name",
		"s.group_id",
		"s.email",
		"s.password",
		"g.name",
	).From("students s")

	query = query.LeftJoin("groups g ON s.group_id = g.id")
	query = query.Where(sq.Eq{"s.uuid": uuid})

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	row := this.db.QueryRow(ctx, sql, args...)

	var result models.Student
	result.Group = &models.Group{}

	if err := row.Scan(
		&result.UUID,
		&result.FirstName,
		&result.LastName,
		&result.MiddleName,
		&result.Group.ID,
		&result.Email,
		&result.Password,
		&result.Group.Name,
	); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	return &result, nil
}

// DeleteStudent implements repositories.StudentRepository.
func (this *studentRepo) DeleteStudent(ctx context.Context, student commands.DeleteStudent) error {
	sql, args, err := sq.Delete("students").Where(sq.Eq{"uuid": student.UUID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	if _, err = this.db.Exec(ctx, sql, args...); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}

// UpdateStudent implements repositories.StudentRepository.
func (this *studentRepo) UpdateStudent(ctx context.Context, student commands.UpdateStudent) error {
	query := sq.Update("students").
		Set("first_name", student.FirstName).
		Set("last_name", student.LastName).
		Set("middle_name", student.MiddleName).
		Set("email", student.Email).
		Set("group_id", student.GroupID).
		Where(sq.Eq{"uuid": student.UUID}).
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

// CreateStudent implements repositories.StudentRepository.
func (this *studentRepo) CreateStudent(ctx context.Context, student commands.CreateStudent) (string, error) {
	sql, args, err := sq.
		Insert("students").
		SetMap(sq.Eq{
			"uuid":        sq.Expr("uuid_generate_v4()"),
			"first_name":  student.FirstName,
			"last_name":   student.LastName,
			"middle_name": student.MiddleName,
			"email":       student.Email,
			"group_id":    student.GroupID,
			"password":    student.Password,
		}).
		Suffix("RETURNING uuid").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}
	fmt.Println("DEBUG", sql, args)

	row := this.db.QueryRow(ctx, sql, args...)

	var id string
	if err := row.Scan(&id); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return id, nil
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
		"s.password",
	).From("students s")

	query = query.LeftJoin("groups g ON s.group_id = g.id")

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
		fmt.Println("here", err)
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
			&stdnt.Password,
		); err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return nil, err
		}

		result = append(result, stdnt)
	}

	return result, nil
}
