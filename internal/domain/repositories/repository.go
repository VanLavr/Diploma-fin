package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/VanLavr/Diploma-fin/utils/config"
)

type Repository interface {
	Connector
	TransactionRepository
	ExamRepository
	StudentRepository
	TeacherRepository
	StudentMailer
	GroupRepository
}

type TransactionRepository interface {
	PerformTransaction(context.Context, func(context.Context) error) error
}

type Connector interface {
	ConnectToPostgres(*config.Config) (*pgxpool.Pool, error)
	CloseConnectionWithPostgres(context.Context) error
}
