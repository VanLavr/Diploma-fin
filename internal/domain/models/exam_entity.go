package models

import "time"

type Exam struct {
	ID   int64
	Name string
}

type Debt struct {
	ID      int64
	Date    time.Time
	Exam    *Exam
	Student *Student
	Teacher *Teacher
}
