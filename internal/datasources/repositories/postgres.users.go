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
}

func (r *postgreUserRepository) CreateUser(ctx context.Context, inDom *V1Domains.UserDomain) (err error) {

	userRecord := records.FromUsersV1Domain(inDom)

	_, err = r.conn.NamedQueryContext(ctx, `INSERT INTO users(id, username, email, password, role_id, created_at) VALUES (uuid_generate_v4(), :username, :email, :password, :role_id, :created_at)`, userRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgreUserRepository) GetUserByEmail(ctx context.Context, inDom *V1Domains.UserDomain) (outDomain V1Domains.UserDomain, err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	err = r.conn.GetContext(ctx, &userRecord, `SELECT * FROM users WHERE "email" = $1`, userRecord.Email)
	if err != nil {
		return V1Domains.UserDomain{}, err
	}

	return userRecord.ToV1Domain(), nil
}

func (r *postgreUserRepository) UpdateUserEmail(ctx context.Context, inDom *V1Domains.UserDomain, newEmail string) (err error) {

	userRecord := records.FromUsersV1Domain(inDom)
	params := map[string]interface{}{"oldEmail": userRecord.Email, "newEmail": newEmail}

	_, err = r.conn.NamedQueryContext(ctx, `UPDATE users SET email = :newEmail WHERE "email" = :oldEmail`, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgreUserRepository) DeleteUser(ctx context.Context, inDom *V1Domains.UserDomain) (err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	_, err = r.conn.NamedQueryContext(ctx, `DELETE FROM users WHERE "id" = :id`, userRecord)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgreUserRepository) GetUserByID(ctx context.Context, inDom *V1Domains.UserDomain) (outDom V1Domains.UserDomain, err error) {
	userRecord := records.FromUsersV1Domain(inDom)

	err = r.conn.GetContext(ctx, &userRecord, `SELECT * FROM users WHERE "id" = $1`, userRecord.Id)
	if err != nil {
		return V1Domains.UserDomain{}, err
	}

	return userRecord.ToV1Domain(), nil
}
