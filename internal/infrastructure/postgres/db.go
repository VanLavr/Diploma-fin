package postgres

import (
	"context"
	"time"

	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/internal/infrastructure/mail"
	"github.com/VanLavr/Diploma-fin/utils/config"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
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

type connector struct {
	pool *pgxpool.Pool
}

func NewConnector(cfg *config.Config) repositories.Connector {
	return &connector{}
}

func (c connector) ConnectToPostgres(cfg *config.Config) (*pgxpool.Pool, error) {
	// Create a configuration from the connection string
	config, err := pgxpool.ParseConfig(cfg.DbString)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// Set connection pool settings (optional)
	config.MaxConns = 10                      // Maximum number of connections in the pool
	config.MinConns = 2                       // Minimum number of connections in the pool
	config.MaxConnLifetime = time.Hour        // Maximum lifetime of a connection
	config.MaxConnIdleTime = time.Minute * 30 // Maximum idle time of a connection

	// Create a connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	c.pool = pool

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	log.Logger.Info("connected to db")
	return pool, nil
}

func (c connector) CloseConnectionWithPostgres(context.Context) error {
	if c.pool != nil {
		c.pool.Close()
		log.Logger.Info("connection to db closed")
	}

	return nil
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
