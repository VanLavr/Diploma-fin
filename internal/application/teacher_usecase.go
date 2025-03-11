package application

import (
	"context"
	"time"

	"github.com/VanLavr/Diploma-fin/internal/application/logic"
	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	valueobjects "github.com/VanLavr/Diploma-fin/internal/domain/value_objects"
	"github.com/VanLavr/Diploma-fin/pkg/errors"
	"github.com/VanLavr/Diploma-fin/pkg/log"
)

type teacherUsecase struct {
	repo repositories.Repository
}

func NewTeacherUsecase() logic.TeacherUsecase {
	return &teacherUsecase{}
}

func (this teacherUsecase) SetDate(ctx context.Context, teacherUUID, date string, debtID int64) error {
	examDate, err := time.Parse(valueobjects.DateLayout, date)
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_APPLICATION, "")
	}

	debts, err := this.repo.GetDebts(ctx, query.GetDebtsFilters{
		DebtIDs: []int64{debtID},
	})
	if err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}
	if len(debts) > 1 || len(debts) == 0 {
		return log.ErrorWrapper(errors.ErroNoItemsFound, errors.ERR_APPLICATION, "")
	}

	if err = this.repo.UpdateDebt(ctx, commands.UpdateDebtByID{
		DebtID:      debtID,
		Date:        examDate,
		TeacherUUID: teacherUUID,
		StudentUUID: debts[0].Student.UUID,
	}); err != nil {
		return log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return nil
}
