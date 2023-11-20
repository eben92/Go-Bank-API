package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	ID        int    `json:"id"`
	Token     string `json:"token"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount   int `json:"to_account"`
	FromAccount int `json:"from_account"`
	Amount      int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    int64     `json:"number"`
	Password  string    `json:"-"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (acc *Account) ValidatePassword(pw string) bool {

	return bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(pw)) == nil

}

func NewAccount(firstName, lastName, password string) (*Account, error) {

	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(100000)),
		Password:  string(encpw),
		CreatedAt: time.Now().UTC(),
	}, nil
}
