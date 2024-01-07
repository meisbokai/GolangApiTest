package requests

import V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"

// Request to create new user
// @Description Request to create new user using username, email and password
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=4" example:"New User 1"`
	Email    string `json:"email" validate:"required,email" example:"newuser1@example.com"`
	Password string `json:"password" validate:"required,min=4" example:"12345"`
} // @name UserCreateRequest

// Mapping Create Request to Domain User
func (user UserCreateRequest) ToV1Domain() *V1Domains.UserDomain {
	return &V1Domains.UserDomain{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		RoleID:   2,
	}
}

// Request to change the email of an existing user
// @Description Request to change password of an existing user
type UserUpdateEmailRequest struct {
	NewEmail string `json:"newEmail" validate:"required,email" example:"test1changed@example.com"`
	OldEmail string `json:"oldEmail" validate:"required,email" example:"test1@example.com"`
} // @name UserUpdateEmailRequest

// Request to login to existing user
// @Description Request to login to existing user using email and password
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"test1@example.com"`
	Password string `json:"password" validate:"required,min=4" example:"12345"`
} // @name UserLoginRequest

// Mapping Login Request to Domain User
func (u *UserLoginRequest) ToV1Domain() *V1Domains.UserDomain {
	return &V1Domains.UserDomain{
		Email:    u.Email,
		Password: u.Password,
	}
}
