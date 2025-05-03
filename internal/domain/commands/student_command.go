package commands

type CreateStudent struct {
	FirstName  string
	LastName   string
	MiddleName string
	GroupID    int64
	Email      string
	Password   string
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

type CreateGroup struct {
	Name string
}

type UpdateGroup struct {
	ID   int64
	Name string
}

type DeleteGroup struct {
	ID int64
}
