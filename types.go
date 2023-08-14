package main

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int64     `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            string    `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int       `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func (a *Account) ValidatePassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            uuid.NewString(),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type TransferRequest struct {
	ToAccount string  `json:"to_account"`
	Amount    float64 `json:"amount"`
}

type LoginRequest struct {
	Number   string `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number string `json:"number"`
	Token  string `json:"token"`
}
