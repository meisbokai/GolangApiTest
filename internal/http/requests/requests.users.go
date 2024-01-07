package requests

import V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"

type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=4" example:"New User 1"`
	Email    string `json:"email" validate:"required,email" example:"newuser1@example.com"`
	Password string `json:"password" validate:"required,min=4" example:"12345"`
}

// Mapping Create Request to Domain User
func (user UserCreateRequest) ToV1Domain() *V1Domains.UserDomain {
	return &V1Domains.UserDomain{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		RoleID:   2,
	}
}

type UserUpdateEmailRequest struct {
	NewEmail string `json:"newEmail" validate:"required,email" example:"test1changed@example.com"`
	OldEmail string `json:"oldEmail" validate:"required,email" example:"test1@example.com"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"test1@example.com"`
	Password string `json:"password" validate:"required,min=4" example:"12345"`
}

// Mapping Login Request to Domain User
func (u *UserLoginRequest) ToV1Domain() *V1Domains.UserDomain {
	return &V1Domains.UserDomain{
		Email:    u.Email,
		Password: u.Password,
	}
}
