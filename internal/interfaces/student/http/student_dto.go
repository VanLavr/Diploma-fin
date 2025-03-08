package http

import "github.com/VanLavr/Diploma-fin/internal/application/exam"

type getAllDebtsDTO struct {
	Err  error  `json:"error"`
	Data []Exam `json:"data"`
}

type Exam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func copyExam(dst *Exam, src *exam.Exam) {
	src.ID = dst.ID
	src.Name = dst.Name
}
