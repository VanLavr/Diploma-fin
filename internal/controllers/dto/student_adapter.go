package dto

import (
	"time"

	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

func DebtDTOFromTypes(src types.Debt) Debt {
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

	var date string
	if src.Date != nil {
		date = src.Date.Format(time.RFC3339)
	}
	return Debt{
		ID:      src.ID,
		Date:    date,
		Exam:    &ex,
		Teacher: &teacher,
		Student: &student,
	}
}

func TypesStudentFromCreateStudentDTO(src CreateStudentDTO) types.Student {
	return types.Student{}
}

func TypeStudentFromUpdateStudentDTO(src UpdateStudentDTO) types.Student {
	return types.Student{}
}

func StudentDTOFromTypes(src types.Student) Student {
	return Student{}
}
