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
	GetAccountByNumber(string) (*Account, error)
	CreateAccount(*Account) error
	DeleteAccount(int) (int, error)
	Transfer(string, float64) (int, error)
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

func (s *PostgresStore) GetAccountByNumber(number string) (*Account, error) {
	rows, err := s.db.Query("select * from accounts where number = $1", number)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %s not found", number)
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		insert into accounts (first_name, last_name, number, encrypted_password, balance, created_at)
		values($1, $2, $3, $4, $5, $6);`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

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

func (s *PostgresStore) Transfer(accountNumber string, amount float64) (int, error) {
	query := `
		update accounts
		set balance = balance + $1
		where number = $2
		returning id;
	`
	rows, err := s.db.Query(query, amount, accountNumber)
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
			encrypted_password varchar(255),
			balance int,
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
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)
	return acc, err
}
