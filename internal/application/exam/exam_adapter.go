package exam

import "github.com/VanLavr/Diploma-fin/internal/domain/exam"

func DomainToApplication(src *exam.Debt, dst *Debt) {
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

	dst.Date = src.Date
}
