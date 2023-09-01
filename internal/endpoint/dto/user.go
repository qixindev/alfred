package dto

type UserDto struct {
	Id               uint   `json:"id"`
	Username         string `json:"username"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	DisplayName      string `json:"displayName"`
	Email            string `json:"email"`
	EmailVerified    bool   `json:"emailVerified"`
	Phone            string `json:"phone"`
	PhoneVerified    bool   `json:"phoneVerified"`
	TwoFactorEnabled bool   `json:"twoFactorEnabled"`
	Avatar           string `json:"avatar"`
	Disabled         bool   `json:"disabled"`
}

type UserProfileDto struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
}

type UserAdminDto struct {
	Id               uint   `json:"id"`
	Username         string `json:"username"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	DisplayName      string `json:"displayName"`
	Email            string `json:"email"`
	EmailVerified    bool   `json:"emailVerified"`
	PasswordHash     string `json:"passwordHash,omitempty"`
	Phone            string `json:"phone"`
	PhoneVerified    bool   `json:"phoneVerified"`
	TwoFactorEnabled bool   `json:"twoFactorEnabled"`
	Disabled         bool   `json:"disabled"`
}
