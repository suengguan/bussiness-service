package models

type Job struct {
	Id          int64
	Name        string
	Description string
	CreatedBy   int64
	CreatedAt   int64
	Project     *Project
	Modules     []*Module
}
