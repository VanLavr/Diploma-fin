package types

import (
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	entities "github.com/VanLavr/Diploma-fin/internal/domain/models"
)

func ExamFromDomain(src *entities.Exam) Exam {
	return Exam{
		ID:   src.ID,
		Name: src.Name,
	}
}

func GroupFromDomain(src *entities.Group) Group {
	return Group{
		ID:   src.ID,
		Name: src.Name,
	}
}

func DebtFromDomain(src *entities.Debt) Debt {
	var (
		ex      = Exam{}
		teacher = Teacher{}
		student = Student{
			Group: &Group{},
		}
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
		if src.Student.Group != nil {
			student.Group.ID = src.Student.Group.ID
			student.Group.Name = src.Student.Group.Name
		}
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
		Password:   src.Password,
	}
}
func StudentFromDomain(src *entities.Student) Student {
	groupName := ""
	var groupID int64
	if src.Group != nil {
		groupName = src.Group.Name
		groupID = src.Group.ID
	}
	return Student{
		UUID:       src.UUID,
		FirstName:  src.FirstName,
		LastName:   src.LastName,
		MiddleName: src.MiddleName,
		Email:      src.Email,
		Password:   src.Password,
		Group: &Group{
			ID:   groupID,
			Name: groupName,
		},
	}
}

func DomainFromStudent(src models.Student) Student {
	var (
		groupID   int64
		groupName string
	)

	if src.Group != nil {
		groupID = src.Group.ID
		groupName = src.Group.Name
	}

	return Student{
		UUID:       src.UUID,
		FirstName:  src.FirstName,
		LastName:   src.LastName,
		MiddleName: src.MiddleName,
		Email:      src.Email,
		Group: &Group{
			ID:   groupID,
			Name: groupName,
		},
	}
}
func DomainFromTeacher(src models.Teacher) Teacher {
	return Teacher{
		UUID:       src.UUID,
		FirstName:  src.FirstName,
		LastName:   src.LastName,
		MiddleName: src.MiddleName,
		Email:      src.Email,
	}
}

func DomainFromExam(src models.Exam) Exam {
	return Exam{
		ID:   src.ID,
		Name: src.Name,
	}
}
