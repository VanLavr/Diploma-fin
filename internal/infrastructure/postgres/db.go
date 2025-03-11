package postgres

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/internal/infrastructure/mail"
	"github.com/VanLavr/Diploma-fin/pkg/config"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	repositories.Connector
	repositories.TransactionRepository
	repositories.ExamRepository
	repositories.StudentRepository
	repositories.TeacherRepository
	repositories.StudentMailer
}

func NewRepository(
	cfg *config.Config,
	transactionRepo repositories.TransactionRepository,
	examRepo repositories.ExamRepository,
	studRepo repositories.StudentRepository,
	teacherRepo repositories.TeacherRepository,
	studMailer repositories.StudentMailer,
) repositories.Repository {
	connector := NewConnector(cfg)
	conn, err := connector.ConnectToPostgres(cfg)
	errors.FatalOnError(err)

	return &repository{
		Connector:             connector,
		TransactionRepository: NewTransaction(conn),
		ExamRepository:        NewExamRepo(conn),
		StudentRepository:     NewStudentRepo(conn),
		TeacherRepository:     NewTeacherRepo(conn),
		StudentMailer:         mail.NewStudentMailer(cfg),
	}
}

type connector struct{}

func NewConnector(cfg *config.Config) repositories.Connector {
	return &connector{}
}

func (c connector) ConnectToPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	panic("not implemented")
}

func (c connector) CloseConnectionWithPostgres(context.Context) error {
	panic("not implemented")
}

type transaction struct {
	db *pgxpool.Pool
}

func NewTransaction(conn *pgxpool.Pool) repositories.TransactionRepository {
	return &transaction{
		db: conn,
	}
}

func (t *transaction) PerformTransaction(ctx context.Context, wrapper func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	defer tx.Rollback(ctx)

	contextWithTransaction := context.WithValue(ctx, valueobjects.TransactionKey{}, tx)

	err = wrapper(contextWithTransaction)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
