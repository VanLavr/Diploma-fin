package http

import (
	"time"

	"github.com/VanLavr/Diploma-fin/internal/application/exam"
)

type getAllDebtsDTO struct {
	Err  error  `json:"error"`
	Data []Debt `json:"data"`
}

type Debt struct {
	ID       int64  `json:"id"`
	ExamID   int64  `json:"exam_id"`
	ExamName string `json:"exam_name"`

	StudentUUID       string `json:"student_uuid"`
	StudentFirstName  string `json:"student_first_name"`
	StudentLastName   string `json:"student_last_name"`
	StudentMiddleName string `json:"student_middle_name"`

	TeacherUUID       string `json:"teacher_uuid"`
	TeacherFirstName  string `json:"teacher_first_name"`
	TeacherLastName   string `json:"teacher_last_name"`
	TeacherMiddleName string `json:"teacher_middle_name"`

	Date string `json:"date"`
}

func copyExam(dst *Debt, src *exam.Debt) {
	dst.ID = src.ID

	dst.ExamID = src.ExamID
	dst.ExamName = src.ExamName

	dst.StudentUUID = src.ExamName
	dst.StudentFirstName = src.StudentFirstName
	dst.StudentLastName = src.StudentLastName
	dst.StudentMiddleName = src.StudentMiddleName

	dst.TeacherUUID = src.TeacherUUID
	dst.TeacherFirstName = src.TeacherFirstName
	dst.TeacherLastName = src.TeacherLastName
	dst.TeacherMiddleName = src.TeacherMiddleName

	dst.Date = src.Date.Format(time.RFC3339)
}
