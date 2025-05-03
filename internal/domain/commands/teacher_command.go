package commands

type CreateTeacher struct {
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Password   string
}

type UpdateTeacher struct {
	UUID       string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
}

type DeleteTeacher struct {
	UUID string
}
