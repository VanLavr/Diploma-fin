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

// CreateTeacher implements logic.TeacherUsecase.
func (t *teacherUsecase) CreateTeacher(ctx context.Context, teacher types.Teacher) (string, error) {
	uuid, err := t.repo.CreateTeacher(ctx, commands.CreateTeacher{
		FirstName:  teacher.FirstName,
		LastName:   teacher.LastName,
		MiddleName: teacher.MiddleName,
		Email:      teacher.Email,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return uuid, nil
}

// DeleteTeacher implements logic.TeacherUsecase.
func (t *teacherUsecase) DeleteTeacher(ctx context.Context, uuid string) error {
	if err := t.repo.DeleteTeacher(ctx, commands.DeleteTeacher{UUID: uuid}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

// GetTeacher implements logic.TeacherUsecase.
func (t *teacherUsecase) GetTeacher(ctx context.Context, uuid string) (types.Teacher, error) {
	teacher, err := t.repo.GetTeacherByUUID(ctx, uuid)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return types.Teacher{}, err
	}

	result := types.TeacherFromDomain(teacher)

	return result, nil
}

// GetTeachers implements logic.TeacherUsecase.
func (t *teacherUsecase) GetTeachers(ctx context.Context, limit int64, offset int64) ([]types.Teacher, error) {
	teachers, err := t.repo.GetTeachers(ctx, query.GetTeachersFilters{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := make([]types.Teacher, len(teachers))
	for i, teacher := range teachers {
		result[i] = types.TeacherFromDomain(&teacher)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// UpdateTeacher implements logic.TeacherUsecase.
func (t *teacherUsecase) UpdateTeacher(ctx context.Context, teacher types.Teacher) error {
	if err := t.repo.UpdateTeacher(ctx, commands.UpdateTeacher{
		UUID:       teacher.UUID,
		FirstName:  teacher.FirstName,
		LastName:   teacher.LastName,
		MiddleName: teacher.MiddleName,
		Email:      teacher.Email,
	}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
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

func (t teacherUsecase) GetTeacherByEmail(ctx context.Context, email string) ([]types.Teacher, error) {
	teachers, err := t.repo.GetTeachers(ctx, query.GetTeachersFilters{
		Emails: []string{email},
	})

	result := make([]types.Teacher, len(teachers))
	for i, teacher := range teachers {
		result[i] = types.TeacherFromDomain(&teacher)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}
