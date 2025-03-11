package dto

type SetDebtDateDTO struct {
	TeacherUUID string `json:"teacher_uuid"`
	DebtID      int64  `json:"debt_id"`
	ExamDate    string `json:"exam_date"`
}

type Exam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
