package dto

type GetStudentDTO struct {
	Err  error   `json:"error"`
	Data Student `json:"data"`
}

type GetAllStudentsDTO struct {
	Err  error     `json:"error"`
	Data []Student `json:"data"`
}

type GetTeacherDTO struct {
	Err  error   `json:"error"`
	Data Teacher `json:"data"`
}

type GetAllTeachersDTO struct {
	Err  error     `json:"error"`
	Data []Teacher `json:"data"`
}

type GetExamDTO struct {
	Err  error `json:"error"`
	Data Exam  `json:"data"`
}

type GetGroupDTO struct {
	Err  error `json:"error"`
	Data Group `json:"data"`
}

type GetDebtDTO struct {
	Err  error `json:"error"`
	Data Debt  `json:"data"`
}

type GetAllExamsDTO struct {
	Err  error  `json:"error"`
	Data []Exam `json:"data"`
}

type GetAllGroupsDTO struct {
	Err  error   `json:"error"`
	Data []Group `json:"data"`
}

type GetAllDebtsDTO struct {
	Err  error  `json:"error"`
	Data []Debt `json:"data"`
}

type CreateStudentResponseDTO struct {
	Err  error  `json:"error"`
	Data string `json:"UUID"`
}

type CreateExamResponseDTO struct {
	Err  error `json:"error"`
	Data int   `json:"id"`
}

type CreateTeacherResponseDTO struct {
	Err  error  `json:"error"`
	Data string `json:"UUID"`
}

type Student struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	GroupName  string `json:"group_name"`
	Email      string `json:"email"`
	GroupID    int64  `json:"group_id"`
}

type Teacher struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
}

type Debt struct {
	ID   int64  `json:"id"`
	Date string `json:"date"`

	Student *Student `json:"student"`
	Teacher *Teacher `json:"teacher"`
	Exam    *Exam    `json:"exam"`
}

type CreateExamDTO struct {
	Name string `json:"name"`
}

type CreateDebtDTO struct {
	ExamID      int64  `json:"id"`
	TeacherUUID string `json:"teacher_uuid"`
	StudentUUID string `json:"student_uuid"`
	Date        string `json:"date"`
}

type CreateStudentDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	GroupID    int64  `json:"group_id"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type CreateTeacherDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateTeacherPasswordDTO struct {
	UUID        string `json:"uuid"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type UpdateStudentPasswordDTO struct {
	UUID        string `json:"uuid"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type UpdateStudentDTO struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	GroupID    int64  `json:"group_id"`
	Email      string `json:"email"`
}

type UpdateTeacherDTO struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
}

type UpdateExamDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateDebtDTO struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"`
	TeacherUUID string `json:"teacher_uuid"`
	StudentUUID string `json:"student_uuid"`
}
