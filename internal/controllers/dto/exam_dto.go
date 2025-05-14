package dto

type SetDebtDateDTO struct {
	TeacherUUID string `json:"teacher_uuid"`
	ExamID      int64  `json:"exam_id"`
	ExamDate    string `json:"exam_date"`
	Address     string `json:"address"`
}

type Exam struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
