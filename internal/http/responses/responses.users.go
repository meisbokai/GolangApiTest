package responses

import (
	"time"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
)

// User information
// @Description Response to a user request containing user information
type UserResponse struct {
	Id        string     `json:"id" example:"1a2b3c"`
	Username  string     `json:"username" example:"John Doe"`
	Email     string     `json:"email" example:"johndoe@example.com"`
	Password  string     `json:"password,omitempty"`
	RoleId    int        `json:"role_id" example:"2"`
	Token     string     `json:"token,omitempty"`
	CreatedAt time.Time  `json:"created_at" example:"2023-12-21 18:58:07.230517+00"`
	UpdatedAt *time.Time `json:"updated_at" example:"2023-12-21 19:20:07.230517+00"`
} // @name UserResponse

func (u *UserResponse) ToV1Domain() V1Domains.UserDomain {
	return V1Domains.UserDomain{
		ID:        u.Id,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		RoleID:    u.RoleId,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromV1Domain(u V1Domains.UserDomain) UserResponse {
	return UserResponse{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Token:     u.Token,
		RoleId:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToResponseList(domains []V1Domains.UserDomain) []UserResponse {
	var result []UserResponse

	for _, val := range domains {
		result = append(result, FromV1Domain(val))
	}

	return result
}
