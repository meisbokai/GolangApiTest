package responses

import (
	"time"

	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
)

type UserResponse struct {
	Id        string
	Username  string
	Email     string
	Password  string
	RoleId    int
	Token     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

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
