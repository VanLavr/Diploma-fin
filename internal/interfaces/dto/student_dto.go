package dto

type GetAllDebtsDTO struct {
	Err  error  `json:"error"`
	Data []Debt `json:"data"`
}

type Student struct {
	UUID       string `json:"uuid"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
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
