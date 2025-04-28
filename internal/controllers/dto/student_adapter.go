package dto

import (
	"time"

	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
)

func ExamDTOFromTypes(src types.Exam) Exam {
	return Exam{
		ID:   src.ID,
		Name: src.Name,
	}
}

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
		teacher.Email = src.Teacher.Email
	}
	if src.Student != nil {
		student.UUID = src.Student.UUID
		student.FirstName = src.Student.FirstName
		student.LastName = src.Student.LastName
		student.MiddleName = src.Student.MiddleName
		student.Email = src.Student.Email
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

func TypesDebtFromCreateDebtDTO(src CreateDebtDTO) (*types.Debt, error) {
	examDate, err := time.Parse(valueobjects.DateLayout, src.Date)
	if err != nil {
		return nil, err
	}

	return &types.Debt{
		Date:    &examDate,
		Exam:    &types.Exam{ID: src.ExamID},
		Student: &types.Student{UUID: src.StudentUUID},
		Teacher: &types.Teacher{UUID: src.TeacherUUID},
	}, nil
}

func TypesExamFromCreateExamDTO(src CreateExamDTO) types.Exam {
	return types.Exam{Name: src.Name}
}

func TypeStudentFromUpdateStudentDTO(src UpdateStudentDTO) types.Student {
	return types.Student{}
}

func StudentDTOFromTypes(src types.Student) Student {
	return Student{}
}

func TeacherDTOFromTypes(src types.Teacher) Teacher {
	return Teacher{}
}

func TypesTeacherFromCreateTeacherDTO(src CreateTeacherDTO) types.Teacher {
	return types.Teacher{}
}

func TypesTeacherFromUpdateTeachertDTO(src UpdateTeacherDTO) types.Teacher {
	return types.Teacher{}
}

func TypesExamFromUpdateExamDTO(src UpdateExamDTO) types.Exam {
	return types.Exam{
		ID:   src.ID,
		Name: src.Name,
	}
}

func TypesDebtFromUpdateDebtDTO(src UpdateDebtDTO) (*types.Debt, error) {
	examDate, err := time.Parse(valueobjects.DateLayout, src.Date)
	if err != nil {
		return nil, err
	}

	return &types.Debt{
		ID:      src.ID,
		Date:    &examDate,
		Student: &types.Student{UUID: src.StudentUUID},
		Teacher: &types.Teacher{UUID: src.TeacherUUID},
	}, nil
}
