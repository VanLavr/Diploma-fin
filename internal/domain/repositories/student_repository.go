package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type StudentRepository interface {
	GetStudents(context.Context, query.GetStudentsFilters) ([]models.Student, error)
}

type StudentMailer interface {
	SendNotification(context.Context, models.Student, string, models.Exam) error
}
