package student

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
)

type StudentRepository interface {
	GetStudents(context.Context, GetStudentsFilters) ([]Student, error)
}

type StudentMailer interface {
	SendNotification(context.Context, Student, string, exam.Exam) error
}
