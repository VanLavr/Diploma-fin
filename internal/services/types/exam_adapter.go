package types

import (
	entities "github.com/VanLavr/Diploma-fin/internal/domain/models"
)

func DebtFromDomain(src *entities.Debt) Debt {
	var (
		ex      = Exam{}
		teacher = Teacher{}
		student = Student{}
	)
	if src.Exam != nil {
		ex.ID = src.Exam.ID
		ex.Name = src.Exam.Name
	}
	if src.Teacher != nil {
		teacher.UUID = src.Teacher.UUID
		teacher.FirstName = src.Teacher.FirstName
		teacher.LastName = src.Teacher.LastName
		teacher.MiddleName = src.Teacher.MiddleName
	}
	if src.Student != nil {
		student.UUID = src.Student.UUID
		student.FirstName = src.Student.FirstName
		student.LastName = src.Student.LastName
		student.MiddleName = src.Student.MiddleName
	}

	return Debt{
		ID:      src.ID,
		Date:    src.Date,
		Exam:    &ex,
		Teacher: &teacher,
		Student: &student,
	}
}
