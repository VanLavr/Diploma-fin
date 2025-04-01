package application

import (
	"context"
	"time"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type teacherUsecase struct {
	repo repositories.Repository
}

func NewTeacherUsecase(repo repositories.Repository) logic.TeacherUsecase {
	return &teacherUsecase{
		repo: repo,
	}
}

func (t teacherUsecase) GetAllDebts(ctx context.Context, UUID string) ([]types.Debt, error) {
	debts, err := t.repo.GetDebts(ctx, query.GetDebtsFilters{
		TeacherUUIDs: []string{UUID},
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

func (this teacherUsecase) SetDate(ctx context.Context, teacherUUID, date string, debtID int64) error {
	examDate, err := time.Parse(valueobjects.DateLayout, date)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return log.ErrorWrapper(err, errors.ERR_APPLICATION, "")
	}

	debts, err := this.repo.GetDebts(ctx, query.GetDebtsFilters{
		DebtIDs: []int64{debtID},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	if len(debts) < 1 || len(debts) == 0 {
		log.Logger.Error("no debts", errors.MethodKey, log.GetMethodName())
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}

	if err = this.repo.UpdateDebt(ctx, commands.UpdateDebtByID{
		DebtID:      debtID,
		Date:        examDate,
		TeacherUUID: teacherUUID,
		StudentUUID: debts[0].Student.UUID,
	}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
