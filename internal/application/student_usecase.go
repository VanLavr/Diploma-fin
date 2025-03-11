package application

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/application/logic"
	"github.com/VanLavr/Diploma-fin/internal/application/types"
	"github.com/VanLavr/Diploma-fin/internal/domain/entities"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type studentUsecase struct {
	examRepo      repositories.ExamRepository
	studentRepo   repositories.StudentRepository
	studentMailer repositories.StudentMailer
	teacherRepo   repositories.TeacherRepository
}

func NewStudentUsecase() logic.StudentUsecase {
	return &studentUsecase{}
}

func (this studentUsecase) GetAllDebts(ctx context.Context, UUID string) ([]types.Debt, error) {
	debts, err := this.examRepo.GetDebts(ctx, query.GetDebtsFilters{
		StudentUUIDs: []string{UUID},
	})

	result := make([]types.Debt, len(debts))
	for i, debt := range debts {
		result[i] = types.DebtFromDomain(&debt)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

func (this studentUsecase) SendNotification(ctx context.Context, UUID string, examID int64) error {
	// get student personal data
	students, err := this.studentRepo.GetStudents(ctx, query.GetStudentsFilters{
		IDs: []string{UUID},
	})
	if len(students) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get debt by id
	exams, err := this.examRepo.GetDebts(ctx, query.GetDebtsFilters{
		ExamIDs: []int64{examID},
	})
	if len(exams) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get teacher personal data
	teachers, err := this.teacherRepo.GetTeachers(ctx, query.GetTeachersFilters{
		UUIDs: []string{exams[0].Teacher.UUID},
	})
	if len(teachers) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// send notification
	err = this.studentMailer.SendNotification(ctx, students[0], teachers[0].Email, entities.Exam{
		ID:   examID,
		Name: exams[0].Exam.Name,
	})
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
