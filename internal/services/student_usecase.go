package application

import (
	"context"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type studentUsecase struct {
	repo repositories.Repository
}

func NewStudentUsecase(repo repositories.Repository) logic.StudentUsecase {
	return &studentUsecase{
		repo: repo,
	}
}

func (this studentUsecase) GetAllDebts(ctx context.Context, UUID string) ([]types.Debt, error) {
	debts, err := this.repo.GetDebts(ctx, query.GetDebtsFilters{
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
	students, err := this.repo.GetStudents(ctx, query.GetStudentsFilters{
		IDs: []string{UUID},
	})
	if len(students) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get debt by id
	exams, err := this.repo.GetDebts(ctx, query.GetDebtsFilters{
		ExamIDs: []int64{examID},
	})
	if len(exams) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// get teacher personal data
	teachers, err := this.repo.GetTeachers(ctx, query.GetTeachersFilters{
		UUIDs: []string{exams[0].Teacher.UUID},
	})
	if len(teachers) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	// send notification
	err = this.repo.SendNotification(ctx, students[0], teachers[0].Email, models.Exam{
		ID:   examID,
		Name: exams[0].Exam.Name,
	})
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
