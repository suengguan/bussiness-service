package models

type Module struct {
	Id          int64
	Name        string
	InputFiles  string
	OutputFiles string
	Parameters  string
	Algorithm   string
	Description string
	Job         *Job
	Pods        []*Pod
}
