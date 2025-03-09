package exam

import "time"

type Debt struct {
	ID int64

	ExamID   int64
	ExamName string

	StudentUUID       string
	StudentFirstName  string
	StudentLastName   string
	StudentMiddleName string

	TeacherUUID       string
	TeacherFirstName  string
	TeacherLastName   string
	TeacherMiddleName string

	Date time.Time
}
