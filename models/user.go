package models

type User struct {
	Id                int64  `json:"id" orm:"column(ID)"`
	Name              string `json:"name" orm:"column(NAME)"`
	Header            string `json:"header" orm:"column(HEADER)"`
	Email             string `json:"email" orm:"column(EMAIL)"`
	Phone             string `json:"phone" orm:"column(PHONE)"`
	Company           string `json:"company" orm:"column(COMPANY)"`
	EncryptedPassword string `json:"encryptedPassword" orm:"column(ENCRYPTED_PASSWORD)"`
	CreatedAt         int64  `json:"createdAt" orm:"column(CREATED_AT)"`
	UpdatedAt         int64  `json:"updatedAt" orm:"column(UPDATED_AT)"`
	Active            bool   `json:"active" orm:"column(ACTIVE)"`
	Role              int    `json:"role" orm:"column(ROLE)"`
}
