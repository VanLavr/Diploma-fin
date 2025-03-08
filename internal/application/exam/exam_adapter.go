package exam

import "github.com/VanLavr/Diploma-fin/internal/domain/exam"

func DomainToApplication(src *exam.Exam, dst *Exam) {
	dst.ID = src.ID
	dst.Name = src.Name
}
