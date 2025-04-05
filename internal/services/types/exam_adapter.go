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

func TeacherFromDomain(src *entities.Teacher) Teacher {
	return Teacher{
		UUID:       src.UUID,
		FirstName:  src.FirstName,
		LastName:   src.LastName,
		MiddleName: src.MiddleName,
		Email:      src.Email,
	}
}
func StudentFromDomain(src *entities.Student) Student {
	return Student{
		UUID:       src.UUID,
		FirstName:  src.FirstName,
		LastName:   src.LastName,
		MiddleName: src.MiddleName,
		Email:      src.Email,
	}
}
