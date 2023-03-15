package models

import "time"

type UserInfo struct {
	Id          uint
	Sub         string
	Name        string
	FirstName   string
	LastName    string
	DisplayName string
	Email       string
	Phone       string
	Picture     string
}

type PhoneVerification struct {
	Id        uint
	Phone     string
	Code      string
	CreatedAt time.Time
}
