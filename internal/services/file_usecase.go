package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"

	"github.com/VanLavr/Diploma-fin/internal/domain/repositories"
	"github.com/VanLavr/Diploma-fin/internal/services/logic"
	"github.com/VanLavr/Diploma-fin/internal/services/types"
	"github.com/VanLavr/Diploma-fin/utils/errors"
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
		if len(parts) < 4 {
			continue
		}

		student := &types.Student{
			LastName:   parts[0],
			FirstName:  parts[1],
			MiddleName: parts[2],
			Group:      &types.Group{Name: parts[3]},
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

	fmt.Println("STUDENTS:")
	for _, stud := range students {
		fmt.Printf("%s, %s, %s, %s, %s || ",
			stud.FirstName,
			stud.LastName,
			stud.MiddleName,
			stud.Email,
			stud.Group.Name,
		)
		fmt.Println()
	}

	fmt.Println("TEACHERS:")
	for _, teac := range teachers {
		fmt.Printf("%s, %s, %s, %s || ",
			teac.FirstName,
			teac.LastName,
			teac.MiddleName,
			teac.Email,
		)
		fmt.Println()
	}

	fmt.Println("DEBTS:")
	for _, debt := range debts {
		fmt.Printf("%s, %s, %s || ",
			debt.Exam.Name,
			debt.Student.FirstName,
			debt.Teacher.FirstName,
		)
		fmt.Println()
	}
	return nil
}

func NewFileUsecase(repo repositories.Repository) logic.FileUsecase {
	return &fileUsecase{
		repo: repo,
	}
}
