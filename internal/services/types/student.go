package types

type Student struct {
	UUID       string
	FirstName  string
	LastName   string
	MiddleName string
	Email      string
	Group      *Group
	Password   string
}

type Group struct {
	ID   int64
	Name string
}
