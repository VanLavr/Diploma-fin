package logic

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

type TeacherUsecase interface {
	SetDate(context.Context, string, string, int64) error
	GetAllDebts(context.Context, string) ([]types.Debt, error)
	GetTeacherByEmail(context.Context, string) ([]types.Teacher, error)
	DeleteTeacher(context.Context, string) error
	UpdateTeacher(context.Context, types.Teacher) error
	CreateTeacher(context.Context, types.Teacher) (string, error)
	GetTeachers(context.Context, int64, int64) ([]types.Teacher, error)
	GetTeacher(context.Context, string) (types.Teacher, error)
	ChangePassword(context.Context, string, string) error
}
