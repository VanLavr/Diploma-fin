package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type StudentUsecase interface {
	GetAllDebts(context.Context, string) ([]types.Debt, error)
	SendNotification(context.Context, string, int64) error
	GetStudentByEmail(context.Context, string) ([]types.Student, error)
	DeleteStudent(context.Context, string) error
	UpdateStudent(context.Context, types.Student) error
	CreateStudent(context.Context, types.Student) (string, error)
	GetStudents(context.Context, int64, int64) ([]types.Student, error)
}
