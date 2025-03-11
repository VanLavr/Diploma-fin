package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type TeacherRepository interface {
	GetTeachers(context.Context, query.GetTeachersFilters) ([]models.Teacher, error)
}
