package main

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    string    `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Number:    uuid.NewString(),
		CreatedAt: time.Now().UTC(),
	}
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TransferRequest struct {
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
}
