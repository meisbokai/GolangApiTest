package records

import (
	"time"
)

type Users struct {
	Id        string     `db:"id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	RoleId    int        `db:"role_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
