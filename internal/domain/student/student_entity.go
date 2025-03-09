package student

type Student struct {
	UUID       string
	FirstName  string
	LastName   string
	MiddleName string
	Group      *Group
}
