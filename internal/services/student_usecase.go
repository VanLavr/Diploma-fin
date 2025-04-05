package application

import (
	"context"
	"fmt"

	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
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
	fmt.Println("10")
	students, err := this.repo.GetStudents(ctx, query.GetStudentsFilters{
		IDs: []string{UUID},
	})
	if len(students) == 0 {
		fmt.Println("11")
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}
	if err != nil {
		fmt.Println("12")
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	fmt.Println("5")

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
	fmt.Println("6")

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
	fmt.Println("7")

	// send notification
	err = this.repo.SendNotification(ctx, students[0], teachers[0].Email, models.Exam{
		ID:   examID,
		Name: exams[0].Exam.Name,
	})
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	fmt.Println("8")

	return nil
}

func (s studentUsecase) GetStudentByEmail(ctx context.Context, email string) ([]types.Student, error) {
	students, err := s.repo.GetStudents(ctx, query.GetStudentsFilters{
		Emails: []string{email},
	})

	result := make([]types.Student, len(students))
	for i, student := range students {
		result[i] = types.StudentFromDomain(&student)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}
