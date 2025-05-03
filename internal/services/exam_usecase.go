package application

import (
	"context"
	"time"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type examUsecase struct {
	repo repositories.Repository
}

// CreateExam implements logic.ExamUsecase.
func (e *examUsecase) CreateExam(ctx context.Context, exam types.Exam) (int64, error) {
	id, err := e.repo.CreateExam(ctx, commands.CreateExam{
		Name: exam.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func NewExamUsecase(repo repositories.Repository) logic.ExamUsecase {
	return &examUsecase{
		repo: repo,
	}
}

// GetDebts implements logic.ExamUsecase.
func (e *examUsecase) GetDebts(ctx context.Context, limit int64, offset int64) ([]types.Debt, error) {
	debts, err := e.repo.GetDebts(ctx, query.GetDebtsFilters{})
	if err != nil {
		return nil, err
	}
	if len(debts) == 0 {
		return nil, errors.ErroNoItemsFound
	}

	result := make([]types.Debt, 0, len(debts))

	for _, debt := range debts {
		var group types.Group
		if debt.Student.Group != nil {
			group.ID = debt.Student.Group.ID
			group.Name = debt.Student.Group.Name
		}
		result = append(result, types.Debt{
			ID:   debt.ID,
			Date: debt.Date,
			Exam: &types.Exam{
				ID:   debt.Exam.ID,
				Name: debt.Exam.Name,
			},
			Student: &types.Student{
				UUID:       debt.Student.UUID,
				FirstName:  debt.Student.FirstName,
				LastName:   debt.Student.LastName,
				MiddleName: debt.Student.MiddleName,
				Email:      debt.Student.Email,
				Group:      &group,
			},
			Teacher: &types.Teacher{
				UUID:       debt.Teacher.UUID,
				FirstName:  debt.Teacher.FirstName,
				LastName:   debt.Teacher.LastName,
				MiddleName: debt.Teacher.MiddleName,
				Email:      debt.Teacher.Email,
			},
		})
	}

	return result, nil
}

func (e *examUsecase) DeleteDebt(ctx context.Context, id int64) error {
	if err := e.repo.DeleteDebt(ctx, commands.DeleteDebt{ID: id}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

func (e *examUsecase) UpdateDebt(ctx context.Context, debt types.Debt) error {
	if debt.Student.UUID != "" {
		students, err := e.repo.GetStudents(ctx, query.GetStudentsFilters{
			IDs: []string{debt.Student.UUID},
		})
		if err != nil {
			return err
		}
		if len(students) == 0 {
			return errors.ErroNoItemsFound
		}

		student := types.DomainFromStudent(students[0])
		debt.Student = &student
	}

	if debt.Teacher.UUID != "" {
		teachers, err := e.repo.GetTeachers(ctx, query.GetTeachersFilters{
			UUIDs: []string{debt.Teacher.UUID},
		})
		if err != nil {
			return err
		}
		if len(teachers) == 0 {
			return errors.ErroNoItemsFound
		}

		teacher := types.DomainFromTeacher(teachers[0])
		debt.Teacher = &teacher
	}

	if debt.Date == nil {
		date := time.Now()
		debt.Date = &date
	}

	err := e.repo.UpdateDebt(ctx, commands.UpdateDebtByID{
		DebtID:      debt.ID,
		Date:        *debt.Date,
		TeacherUUID: debt.Teacher.UUID,
		StudentUUID: debt.Student.UUID,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

func (e *examUsecase) CreateDebt(ctx context.Context, debt types.Debt) (int64, error) {
	if debt.Exam.ID != 0 {
		exam, err := e.repo.GetExamByID(ctx, query.GetExamsFilters{
			IDs: []int64{debt.Exam.ID},
		})
		if err != nil {
			return 0, err
		}

		exm := types.DomainFromExam(*exam)
		debt.Exam = &exm
	}

	if debt.Student.UUID != "" {
		students, err := e.repo.GetStudents(ctx, query.GetStudentsFilters{
			IDs: []string{debt.Student.UUID},
		})
		if err != nil {
			return 0, err
		}
		if len(students) == 0 {
			return 0, errors.ErroNoItemsFound
		}

		student := types.DomainFromStudent(students[0])
		debt.Student = &student
	}

	if debt.Teacher.UUID != "" {
		teachers, err := e.repo.GetTeachers(ctx, query.GetTeachersFilters{
			UUIDs: []string{debt.Teacher.UUID},
		})
		if err != nil {
			return 0, err
		}
		if len(teachers) == 0 {
			return 0, errors.ErroNoItemsFound
		}

		teacher := types.DomainFromTeacher(teachers[0])
		debt.Teacher = &teacher
	}

	if debt.Date == nil {
		date := time.Now()
		debt.Date = &date
	}

	id, err := e.repo.CreateDebt(ctx, commands.CreateDebt{
		ExamID:      debt.Exam.ID,
		StudentUUID: debt.Student.UUID,
		TeacherUUID: debt.Teacher.UUID,
		Date:        *debt.Date,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func (e *examUsecase) GetDebt(ctx context.Context, id int64) (*types.Debt, error) {
	debts, err := e.repo.GetDebts(ctx, query.GetDebtsFilters{DebtIDs: []int64{id}})
	if err != nil {
		return nil, err
	}
	if len(debts) == 0 {
		return nil, errors.ErroNoItemsFound
	}

	var group types.Group
	if debts[0].Student.Group != nil {
		group.ID = debts[0].Student.Group.ID
		group.Name = debts[0].Student.Group.Name
	}

	return &types.Debt{
		ID:   id,
		Date: debts[0].Date,
		Exam: &types.Exam{
			ID:   debts[0].Exam.ID,
			Name: debts[0].Exam.Name,
		},
		Student: &types.Student{
			UUID:       debts[0].Student.UUID,
			FirstName:  debts[0].Student.FirstName,
			LastName:   debts[0].Student.LastName,
			MiddleName: debts[0].Student.LastName,
			Email:      debts[0].Student.Email,
			Group:      &group,
		},
		Teacher: &types.Teacher{
			UUID:       debts[0].Teacher.UUID,
			FirstName:  debts[0].Teacher.FirstName,
			LastName:   debts[0].Teacher.LastName,
			MiddleName: debts[0].Teacher.MiddleName,
			Email:      debts[0].Teacher.Email,
		},
	}, nil
}

// DeleteExam implements logic.ExamUsecase.
func (e *examUsecase) DeleteExam(ctx context.Context, id int64) error {
	if err := e.repo.DeleteExam(ctx, commands.DeleteExam{ID: id}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

// GetExam implements logic.ExamUsecase.
func (e *examUsecase) GetExam(ctx context.Context, id int64) (*types.Exam, error) {
	exam, err := e.repo.GetExamByID(ctx, query.GetExamsFilters{
		IDs: []int64{id},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := types.ExamFromDomain(exam)

	return &result, nil
}

// GetExams implements logic.ExamUsecase.
func (e *examUsecase) GetExams(ctx context.Context, limit int64, offset int64) ([]types.Exam, error) {
	exams, err := e.repo.GetExams(ctx, query.GetExamsFilters{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := make([]types.Exam, len(exams))
	for i, exam := range exams {
		result[i] = types.ExamFromDomain(&exam)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// UpdateExam implements logic.ExamUsecase.
func (e *examUsecase) UpdateExam(ctx context.Context, exam types.Exam) error {
	err := e.repo.UpdateExam(ctx, commands.UpdateExamByID{ID: exam.ID, Name: exam.Name})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}
