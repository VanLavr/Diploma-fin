package student

import (
	"context"

	ex "github.com/VanLavr/Diploma-fin/internal/application/exam"
	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
	"github.com/VanLavr/Diploma-fin/internal/domain/student"
	"github.com/VanLavr/Diploma-fin/internal/domain/teacher"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
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
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	data := make([]ex.Debt, len(debts))
	for i, debt := range debts {
		ex.DomainToApplication(&debt, &data[i])
	}

	return data, nil
}

func (this usecase) SendNotification(ctx context.Context, UUID string, examID int64) error {
	// get student personal data
	students, err := this.studentRepo.GetStudents(ctx, student.GetStudentsFilters{
		IDs: []string{UUID},
	})
	if len(students) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get debt by id
	exams, err := this.examRepo.GetDebts(ctx, exam.GetDebtsFilters{
		ExamIDs: []int64{examID},
	})
	if len(exams) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get teacher personal data
	teachers, err := this.teacherRepo.GetTeachers(ctx, teacher.GetTeachersFilters{
		UUIDs: []string{exams[0].TeacherUUID},
	})
	if len(teachers) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// send notification
	err = this.studentMailer.SendNotification(ctx, students[0], teachers[0].Email, exam.Exam{
		ID:   examID,
		Name: exams[0].ExamName,
	})
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
