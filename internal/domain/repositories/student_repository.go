package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/entities"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type StudentRepository interface {
	GetStudents(context.Context, query.GetStudentsFilters) ([]entities.Student, error)
}

type StudentMailer interface {
	SendNotification(context.Context, entities.Student, string, entities.Exam) error
}
