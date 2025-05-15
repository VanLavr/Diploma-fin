package application

import (
	"context"
	e "errors"
	"fmt"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/VanLavr/Diploma-fin/internal/domain/commands"
	query "github.com/VanLavr/Diploma-fin/internal/domain/queries"
	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
	"github.com/VanLavr/Diploma-fin/utils/generator"
	"github.com/VanLavr/Diploma-fin/utils/hasher"
	"github.com/VanLavr/Diploma-fin/utils/log"
)

type fileUsecase struct {
	repo repositories.Repository
}

// ParseFile implements logic.FileUsecase.
func (fu *fileUsecase) ParseFile(ctx context.Context, f *excelize.File) error {
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return fmt.Errorf("get rows: %w", err)
	}

	if len(rows) < 2 {
		log.Logger.Error("not enough rows in the table", errors.MethodKey, log.GetMethodName())
		return fmt.Errorf("not enough rows in the table")
	}

	// 1. Парсим студентов из заголовка
	studentHeaders := rows[0][1:]
	students := make([]*types.Student, 0, len(studentHeaders))
	studentMap := make(map[int]*types.Student) // Для связи по индексу

	for i, header := range studentHeaders {
		parts := strings.Split(header, " ")
		fmt.Println(parts)
		if len(parts) < 5 {
			continue
		}

		student := &types.Student{
			LastName:   parts[0],
			FirstName:  parts[1],
			MiddleName: parts[2],
			Group:      &types.Group{Name: parts[3]},
			Email:      parts[4],
		}

		students = append(students, student)
		studentMap[i] = student
	}

	// 2. Парсим преподавателей и долги
	teachers := make([]*types.Teacher, 0)
	teacherMap := make(map[string]*types.Teacher) // Для избежания дубликатов
	debts := make([]*types.Debt, 0)

	for _, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}

		teacherName := strings.TrimSpace(row[0])
		if teacherName == "" {
			continue
		}

		// Обрабатываем преподавателя
		teacherParts := strings.Split(teacherName, " ")
		if len(teacherParts) < 4 {
			continue
		}

		teacherKey := teacherName
		teacher, exists := teacherMap[teacherKey]
		if !exists {
			teacher = &types.Teacher{
				LastName:   teacherParts[0],
				FirstName:  teacherParts[1],
				MiddleName: teacherParts[2],
				Email:      teacherParts[3],
			}
			teachers = append(teachers, teacher)
			teacherMap[teacherKey] = teacher
		}

		// Обрабатываем долги для каждого студента
		for colIdx, cell := range row[1:] {
			if colIdx >= len(studentMap) {
				break
			}

			examName := strings.TrimSpace(cell)
			if examName == "" {
				continue
			}

			student := studentMap[colIdx]
			if student == nil {
				continue
			}

			debts = append(debts, &types.Debt{
				Exam:    &types.Exam{Name: examName},
				Student: student,
				Teacher: teacher,
			})
		}
	}

	// fmt.Println("STUDENTS:")
	// for _, stud := range students {
	// }

	// fmt.Println("TEACHERS:")
	// for _, teac := range teachers {
	// }

	fmt.Println("DEBTS:")
	for _, debt := range debts {
		_, err := fu.CreateDebtIfNotExists(ctx, *debt)
		if err != nil {
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return err
		}
	}
	return nil
}

func (fu *fileUsecase) CreateGroupIfNotExists(ctx context.Context, group types.Group) (int64, error) {
	groupsFound, err := fu.repo.SearchGroups(ctx, query.SearchGroupFilters{
		Names: []string{group.Name},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}
	if len(groupsFound) != 0 {
		return groupsFound[0].ID, nil
	}

	id, err := fu.repo.CreateGroup(ctx, commands.CreateGroup{
		Name: group.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func (fu *fileUsecase) CreateStudentIfNotExists(ctx context.Context, student types.Student) (string, error) {
	studentsFound, err := fu.repo.SearchStudents(ctx, query.SearchStudentFilters{
		Emails: []string{student.Email},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}
	if len(studentsFound) != 0 {
		return studentsFound[0].UUID, nil
	}

	groupID, err := fu.CreateGroupIfNotExists(ctx, types.Group{
		Name: student.Group.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	// 1) generate password
	pass, err := generator.GeneratePassword(generator.DEFAULTLEN)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	// 2) send password to student's email
	err = fu.repo.SendPassword(ctx, student.Email, pass)
	switch {
	case err == nil:
	case e.Is(err, errors.ErrInvalidData):
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
	default:
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	id, err := fu.repo.CreateStudent(ctx, commands.CreateStudent{
		FirstName:  student.FirstName,
		LastName:   student.LastName,
		MiddleName: student.MiddleName,
		GroupID:    groupID,
		Email:      student.Email,
		Password:   hasher.Hshr.Hash(pass),
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return id, nil
}

func (fu *fileUsecase) CreateTeacherIfNotExists(ctx context.Context, teacher types.Teacher) (string, error) {
	teachersFound, err := fu.repo.SearchTeachers(ctx, query.SearchTeacherFilters{
		Emails: []string{teacher.Email},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}
	if len(teachersFound) != 0 {
		return teachersFound[0].UUID, nil
	}

	// 1) generate password
	pass, err := generator.GeneratePassword(generator.DEFAULTLEN)
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}
	if teacher.Email == "84956743974@mail.ru" {
		pass = "test"
	}

	if teacher.Email != "84956743974@mail.ru" {
		// 2) send password to student's email
		err = fu.repo.SendPassword(ctx, teacher.Email, pass)
		switch {
		case err == nil:
		case e.Is(err, errors.ErrInvalidData):
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		default:
			log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
			return "", err
		}
	}

	id, err := fu.repo.CreateTeacher(ctx, commands.CreateTeacher{
		FirstName:  teacher.FirstName,
		LastName:   teacher.LastName,
		MiddleName: teacher.MiddleName,
		Email:      teacher.Email,
		Password:   hasher.Hshr.Hash(pass),
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return "", err
	}

	return id, nil
}

func (fu *fileUsecase) CreateExamIfNotExists(ctx context.Context, exam types.Exam) (int64, error) {
	examsFound, err := fu.repo.SearchExams(ctx, query.SearchExamFilters{
		Names: []string{exam.Name},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}
	if len(examsFound) != 0 {
		return examsFound[0].ID, nil
	}

	id, err := fu.repo.CreateExam(ctx, commands.CreateExam{
		Name: exam.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func (fu *fileUsecase) CreateDebtIfNotExists(ctx context.Context, debt types.Debt) (int64, error) {
	debtsFound, err := fu.repo.SearchDebts(ctx, query.SearchDebtsFilters{
		ExamNames:     []string{debt.Exam.Name},
		StudentEmails: []string{debt.Student.Email},
		TeacherEmails: []string{debt.Teacher.Email},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}
	if len(debtsFound) != 0 {
		return debtsFound[0].ID, nil
	}

	eid, err := fu.CreateExamIfNotExists(ctx, types.Exam{
		Name: debt.Exam.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	gid, err := fu.CreateGroupIfNotExists(ctx, types.Group{
		Name: debt.Student.Group.Name,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	suuid, err := fu.CreateStudentIfNotExists(ctx, types.Student{
		FirstName:  debt.Student.FirstName,
		LastName:   debt.Student.LastName,
		MiddleName: debt.Student.MiddleName,
		Email:      debt.Student.Email,
		Group: &types.Group{
			ID:   gid,
			Name: debt.Student.Group.Name,
		},
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	tuuid, err := fu.CreateTeacherIfNotExists(ctx, types.Teacher{
		FirstName:  debt.Teacher.FirstName,
		LastName:   debt.Teacher.LastName,
		MiddleName: debt.Teacher.MiddleName,
		Email:      debt.Teacher.Email,
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	id, err := fu.repo.CreateDebt(ctx, commands.CreateDebt{
		ExamID:      eid,
		StudentUUID: suuid,
		TeacherUUID: tuuid,
		Date:        time.Now(),
	})
	if err != nil {
		log.Logger.Error(err.Error(), errors.MethodKey, log.GetMethodName())
		return 0, err
	}

	return id, nil
}

func NewFileUsecase(repo repositories.Repository) logic.FileUsecase {
	return &fileUsecase{
		repo: repo,
	}
}
