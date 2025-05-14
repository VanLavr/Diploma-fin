package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type StudentRepository interface {
	GetStudentByUUID(context.Context, string) (*models.Student, error)
	GetStudents(context.Context, query.GetStudentsFilters) ([]models.Student, error)
	CreateStudent(context.Context, commands.CreateStudent) (string, error)
	UpdateStudent(context.Context, commands.UpdateStudent) error
	DeleteStudent(context.Context, commands.DeleteStudent) error
	SearchStudents(context.Context, query.SearchStudentFilters) ([]models.Student, error)
	ChangeStudentPassword(context.Context, string, string) error
	GetAmountOfDebtsForStudent(context.Context, string) (int64, error)
}

type StudentMailer interface {
	SendNotification(context.Context, models.Student, string, models.Exam) error
	SendPassword(context.Context, string, string) error
}
