package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type TeacherRepository interface {
	GetTeachers(context.Context, query.GetTeachersFilters) ([]models.Teacher, error)
	GetTeacherByUUID(context.Context, string) (*models.Teacher, error)
	CreateTeacher(context.Context, commands.CreateTeacher) (string, error)
	UpdateTeacher(context.Context, commands.UpdateTeacher) error
	DeleteTeacher(context.Context, commands.DeleteTeacher) error
	SearchTeachers(context.Context, query.SearchTeacherFilters) ([]models.Teacher, error)
	ChangeTeacherPassword(context.Context, string, string) error
}
