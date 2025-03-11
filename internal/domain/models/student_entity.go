package models

type Student struct {
	UUID       string
	FirstName  string
	LastName   string
	MiddleName string
	Group      *Group
}

type Group struct {
	ID   int64
	Name string
}
