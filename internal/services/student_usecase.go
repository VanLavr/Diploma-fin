package application

import (
	"context"
	"fmt"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	"github.com/VanLavr/Diploma-fin/internal/domain/models"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/hasher"
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

// GetAmountOfDebts implements logic.StudentUsecase.
func (this *studentUsecase) GetAmountOfDebts(ctx context.Context, uuid string) (int64, error) {
	return this.repo.GetAmountOfDebtsForStudent(ctx, uuid)
}

// ChangePassword implements logic.StudentUsecase.
func (this *studentUsecase) ChangePassword(ctx context.Context, uuid, password string) error {
	return this.repo.ChangeStudentPassword(ctx, uuid, password)
}

// CreateStudent implements logic.StudentUsecase.
func (this *studentUsecase) CreateStudent(ctx context.Context, student types.Student) (string, error) {
	uuid, err := this.repo.CreateStudent(ctx, commands.CreateStudent{
		FirstName:  student.FirstName,
		LastName:   student.LastName,
		MiddleName: student.MiddleName,
		GroupID:    student.Group.ID,
		Email:      student.Email,
		Password:   hasher.Hshr.Hash(student.Password),
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return uuid, nil
}

// DeleteStudent implements logic.StudentUsecase.
func (this *studentUsecase) DeleteStudent(ctx context.Context, uuid string) error {
	if err := this.repo.DeleteStudent(ctx, commands.DeleteStudent{UUID: uuid}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
}

// GetStudent implements logic.StudentUsecase.
func (this *studentUsecase) GetStudent(ctx context.Context, uuid string) (*types.Student, error) {
	student, err := this.repo.GetStudentByUUID(ctx, uuid)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := types.StudentFromDomain(student)

	return &result, nil
}

// GetStudents implements logic.StudentUsecase.
func (this *studentUsecase) GetStudents(ctx context.Context, limit int64, offset int64) ([]types.Student, error) {
	students, err := this.repo.GetStudents(ctx, query.GetStudentsFilters{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return nil, err
	}

	result := make([]types.Student, len(students))
	for i, student := range students {
		result[i] = types.StudentFromDomain(&student)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}

// UpdateStudent implements logic.StudentUsecase.
func (this *studentUsecase) UpdateStudent(ctx context.Context, student types.Student) error {
	if err := this.repo.UpdateStudent(ctx, commands.UpdateStudent{
		UUID:       student.UUID,
		FirstName:  student.FirstName,
		LastName:   student.LastName,
		MiddleName: student.MiddleName,
		GroupID:    student.Group.ID,
		Email:      student.Email,
	}); err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return err
	}

	return nil
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

	for i, debt := range result {
		debtsWithGroupsByExam, err := this.repo.GetDebts(ctx, query.GetDebtsFilters{
			ExamIDs: []int64{debt.Exam.ID},
		})
		if err != nil {
			return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
		}

		uniqueGroups := make(map[int64]models.Group)
		for _, debt := range debtsWithGroupsByExam {
			uniqueGroups[debt.Student.Group.ID] = models.Group{
				ID:   debt.Student.Group.ID,
				Name: debt.Student.Group.Name,
			}
		}
		groupList := make([]types.Group, 0, len(uniqueGroups))
		for _, group := range uniqueGroups {
			groupList = append(groupList, types.Group{
				ID:   group.ID,
				Name: group.Name,
			})
		}

		result[i].Groups = groupList
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
	fmt.Println("got stud", students)

	result := make([]types.Student, len(students))
	for i, student := range students {
		result[i] = types.StudentFromDomain(&student)
	}

	if err != nil {
		return nil, log.ErrorWrapper(err, errors.ERR_INFRASTRUCTURE, "")
	}

	return result, nil
}
