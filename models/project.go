package models

type Project struct {
	Id          int64
	Name        string
	CreatedBy   int64
	CreatedAt   int64
	Description string
	User        *User
	Jobs        []*Job
}
