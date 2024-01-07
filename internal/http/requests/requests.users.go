package requests

import V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"

type UserCreateRequest struct {
	Username string
	Email    string
	Password string
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
