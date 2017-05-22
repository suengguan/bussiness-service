package models

type Response struct {
	Status     int
	RetryCount int
	Reason     string
	Result     string
}
