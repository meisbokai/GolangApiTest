package v1

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/meisbokai/GolangApiTest/internal/datasources/records"
	V1Domains "github.com/meisbokai/GolangApiTest/internal/http/domains/v1"
)

type postgreUserRepository struct {
	conn *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) V1Domains.UserRepository {
	return &postgreUserRepository{
		conn: conn,
	}
}

func (r *postgreUserRepository) GetAllUsers(ctx context.Context) (outDom []V1Domains.UserDomain, err error) {

	userRecords := []records.Users{}

	err = r.conn.SelectContext(ctx, &userRecords, `SELECT * FROM users`)
	if err != nil {
		return []V1Domains.UserDomain{}, err
	}

	return records.ToArrayOfUsersV1Domain(&userRecords), nil
func (r *postgreUserRepository) GetUserByEmail(ctx context.Context, inDom *V1Domains.UserDomain) (outDomain V1Domains.UserDomain, err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	err = r.conn.GetContext(ctx, &userRecord, `SELECT * FROM users WHERE "email" = $1`, userRecord.Email)
	if err != nil {
		return V1Domains.UserDomain{}, err
	}

	return userRecord.ToV1Domain(), nil
}
