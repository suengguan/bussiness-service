package models

type PodStatus struct {
	Id       int64  `json:"id"`
	Status   int    `json:"status"` // run, stop, pause,finish
	Progress int    `json:"progress"`
	Reason   string `json:"reason"`
}
