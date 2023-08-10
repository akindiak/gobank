package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(int) (int, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	godotenv.Load(".env")
	connStr := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select * from accounts where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)

	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		insert into accounts (first_name, last_name, number, balance, created_at)
		values($1, $2, $3, $4, $5);`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) (int, error) {
	rows, err := s.db.Query("delete from accounts where id = $1 returning id", id)
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		return id, err
	}
	return 0, err
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
		create table if not exists accounts (
			id serial not null primary key,
			first_name varchar(255),
			last_name varchar(255),
			number varchar(255) not null,
			balance double precision,
			created_at timestamp
		);`

	_, err := s.db.Exec(query)
	return err
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := &Account{}
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Balance,
		&acc.CreatedAt,
	)
	return acc, err
}
