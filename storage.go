package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
	DropDB() error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account 
	(first_name, last_name, number, password, balance, created_at)
	values
	($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Password, acc.Balance, acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {

	_, err := s.GetAccountByID(id)

	if err != nil {
		return err
	}

	_, qError := s.db.Query("DELETE FROM account where id = $1", id)

	if qError != nil {
		return fmt.Errorf("account with id %d not found", id)
	}

	return nil
}

func (s *PostgresStore) DropDB() error {
	_, err := s.db.Query("DROP TABLE account")
	if err != nil {
		return fmt.Errorf("account with id not found")
	}

	fmt.Println("table dropped")

	return nil
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account WHERE number = $1", number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account with number [%d] not found", number)

}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account where id = $1", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{} // same as new([]Account)

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}

	return accounts, nil
}

func (s *PostgresStore) Init() error {

	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		password varchar(100),
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)

	return err

}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Password,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
