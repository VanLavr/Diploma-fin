package repositories

import "context"

type Repository interface {
	Connector
	TransactionRepository
	ExamRepository
	StudentRepository
	TeacherRepository
	StudentMailer
}

type TransactionRepository interface {
	PerformTransaction(context.Context, func(context.Context) error) error
}
