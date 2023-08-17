package model

import (
	"accounts/internal/endpoint/dto"
)

type User struct {
	Id               uint   `gorm:"primaryKey;autoIncrement;not null" json:"id"`
	Username         string `json:"username"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	DisplayName      string `json:"displayName" binding:"required"`
	Email            string `json:"email"`
	EmailVerified    bool   `json:"emailVerified"`
	PasswordHash     string `json:"passwordHash,omitempty"`
	Phone            string `json:"phone"`
	PhoneVerified    bool   `json:"phoneVerified"`
	TwoFactorEnabled bool   `json:"twoFactorEnabled"`
	Disabled         bool   `json:"disabled"`

	TenantId uint `gorm:"primaryKey"`
	Tenant   Tenant
	Role     string
}

func (u *User) Name() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

func (u *User) Dto() dto.UserDto {
	return dto.UserDto{
		Id:               u.Id,
		Username:         u.Username,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		DisplayName:      u.DisplayName,
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
		Username:    u.Username,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		Phone:       u.Phone,
	}
}

func (u *User) AdminDto() dto.UserAdminDto {
	return dto.UserAdminDto{
		Id:               u.Id,
		Username:         u.Username,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		DisplayName:      u.DisplayName,
		Email:            u.Email,
		EmailVerified:    u.EmailVerified,
		Phone:            u.Phone,
		PhoneVerified:    u.PhoneVerified,
		TwoFactorEnabled: u.TwoFactorEnabled,
		Disabled:         u.Disabled,
	}
}

func User2AdminDto(u User) dto.UserAdminDto {
	return u.AdminDto()
}
