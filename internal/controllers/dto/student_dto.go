package dto

type GetAllStudentsDTO struct {
	Err  error     `json:"error"`
	Data []Student `json:"data"`
}

type GetAllDebtsDTO struct {
	Err  error  `json:"error"`
	Data []Debt `json:"data"`
}

type CreateStudentResponseDTO struct {
	Err  error  `json:"error"`
	Data string `json:"UUID"`
}

type Student struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	GroupName  string `json:"group_name"`
}

type Teacher struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
}

type Debt struct {
	ID   int64  `json:"id"`
	Date string `json:"date"`

	Student *Student `json:"student"`
	Teacher *Teacher `json:"teacher"`
	Exam    *Exam    `json:"exam"`
}

type CreateStudentDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Group      Group  `json:"group"`
	Email      string `json:"email"`
}

type Group struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UpdateStudentDTO struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	Group      Group  `json:"group"`
	Email      string `json:"email"`
}
