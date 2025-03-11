package repositories

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/entities"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
)

type TeacherRepository interface {
	GetTeachers(context.Context, query.GetTeachersFilters) ([]entities.Teacher, error)
}
