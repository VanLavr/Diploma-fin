package teacher

import "context"

type TeacherRepository interface {
	GetTeachers(context.Context, GetTeachersFilters) ([]Teacher, error)
}
