package student

import (
	"context"

	ex "github.com/VanLavr/Diploma-fin/internal/application/exam"
	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
	"github.com/VanLavr/Diploma-fin/internal/domain/student"
	"github.com/VanLavr/Diploma-fin/internal/domain/teacher"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
)

type StudentUsecase interface {
	GetAllDebts(context.Context, string) ([]ex.Debt, error)
	SendNotification(context.Context, string, int64) error
}

type usecase struct {
	examRepo      exam.ExamRepository
	studentRepo   student.StudentRepository
	studentMailer student.StudentMailer
	teacherRepo   teacher.TeacherRepository
}

func NewStudentUsecase() StudentUsecase {
	return &usecase{}
}

func (this usecase) GetAllDebts(ctx context.Context, UUID string) ([]ex.Debt, error) {
	debts, err := this.examRepo.GetDebts(ctx, exam.GetDebtsFilters{
		StudentUUIDs: []string{UUID},
	})
	if err != nil {
		return nil, err
	}

	data := make([]ex.Debt, len(debts))
	for i, debt := range debts {
		ex.DomainToApplication(&debt, &data[i])
	}

	return data, nil
}

func (this usecase) SendNotification(ctx context.Context, UUID string, examID int64) (err error) {
	// get student personal data
	students, err := this.studentRepo.GetStudents(ctx, student.GetStudentsFilters{
		IDs: []string{UUID},
	})
	if len(students) == 0 {
		return errors.ErroNoItemsFound
	}
	if err != nil {
		return
	}

	// get debt by id
	exams, err := this.examRepo.GetDebts(ctx, exam.GetDebtsFilters{
		ExamIDs: []int64{examID},
	})
	if len(exams) == 0 {
		return errors.ErroNoItemsFound
	}
	if err != nil {
		return
	}

	// get teacher personal data
	teachers, err := this.teacherRepo.GetTeachers(ctx, teacher.GetTeachersFilters{
		UUIDs: []string{exams[0].TeacherUUID},
	})
	if len(teachers) == 0 {
		return errors.ErroNoItemsFound
	}
	if err != nil {
		return
	}

	// send notification
	return
}
