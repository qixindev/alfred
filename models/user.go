package models

import "accounts/models/dto"

type User struct {
	Id               uint
	Username         string
	FirstName        string
	LastName         string
	DisplayName      string
	Email            string
	EmailVerified    bool
	PasswordHash     string
	Phone            string
	PhoneVerified    bool
	TwoFactorEnabled bool
	Disabled         bool

	TenantId uint
	Tenant   Tenant
}

func (u *User) Dto() dto.UserDto {
	return dto.UserDto{
		Id:               u.Id,
		Username:         u.Username,
		Email:            u.Email,
		EmailVerified:    u.EmailVerified,
		Phone:            u.Phone,
		PhoneVerified:    u.PhoneVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		Disabled:         u.Disabled,
	}
}

func (u *User) ProfileDto() dto.UserProfileDto {
	return dto.UserProfileDto{
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
	}
}
