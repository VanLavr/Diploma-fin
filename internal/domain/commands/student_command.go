package commands

type CreateStudent struct {
	FirstName  string
	LastName   string
	MiddleName string
	GroupID    int64
	Email      string
}

type UpdateStudent struct {
	UUID       string
	FirstName  string
	LastName   string
	MiddleName string
	GroupID    int64
	Email      string
}

type DeleteStudent struct {
	UUID string
}
