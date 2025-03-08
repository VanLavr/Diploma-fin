package student

import (
	"context"

	ex "github.com/VanLavr/Diploma-fin/internal/application/exam"
	"github.com/VanLavr/Diploma-fin/internal/domain/exam"
)

type StudentUsecase interface {
	GetAllDebts(context.Context, string) ([]ex.Exam, error)
}

type usecase struct {
	examRepo exam.ExamRepository
}

func NewStudentUsecase() StudentUsecase {
	return &usecase{}
}

func (this usecase) GetAllDebts(ctx context.Context, UUID string) ([]ex.Exam, error) {
	debts, err := this.examRepo.GetAllDebts(ctx, exam.GetAllDebtsFilter{
		StudentUUID: UUID,
	})
	if err != nil {
		return nil, err
	}

	data := make([]ex.Exam, len(debts))
	for i, debt := range debts {
		ex.DomainToApplication(&debt, &data[i])
	}

	return data, nil
}
