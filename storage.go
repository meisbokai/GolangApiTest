package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUser(*User) error
	DeleteUser(int) error
	UpdateUser(*User) error
	GetUserByID(int) (*User, error)
	GetUsers() ([]*User, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=golangapitest sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Postgres connection established")

	return &PostgresStorage{db: db}, nil

}

func (s *PostgresStorage) CreateUser(user *User) error {
	query := `
	insert into users 
	(username, email, created_at)
	values
	($1, $2, $3)
	`

	resp, err := s.db.Query(query,
		user.Username,
		user.Email,
		user.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStorage) UpdateUser(*User) error {
	return nil
}
func (s *PostgresStorage) DeleteUser(id int) error {
	return nil
}
func (s *PostgresStorage) GetUserByID(id int) (*User, error) {
	// Assume that ID is unique since it is a primary key in the db.
	// Query for a single row
	row := s.db.QueryRow(`select * from users where id = $1`, id)

	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("User with ID %d not found", id)
	}

	return user, nil
}

func (s *PostgresStorage) GetUsers() ([]*User, error) {
	rows, err := s.db.Query(`SELECT * from users`)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for rows.Next() {
		user := new(User)
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	return users, nil
}

func (s *PostgresStorage) Init() error {
	return s.createUserTable()
}

func (s *PostgresStorage) createUserTable() error {
	// Note: names that are not all lowercase need to be double quoted
	// Note: 'user' is a reserved word in postgres
	query := `create table if not exists users (
		id SERIAL PRIMARY KEY,
        username VARCHAR(50),
        email VARCHAR(50),
        created_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	if err == nil {
		log.Println("Postgres Storage Initialized.")
	}

	return err
}

// docker run --name some-postgres -e POSTGRES_PASSWORD=golangapitest -p 5432:5432 -d postgres
