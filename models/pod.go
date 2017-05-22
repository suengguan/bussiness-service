package models

type Pod struct {
	Id        int64
	Name      string
	PodName   string
	RcName    string
	SvcName   string
	DataRange string
	Cpu       float64
	Memory    float64
	Module    *Module
}
