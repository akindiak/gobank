package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHandleFunc(s.handleAccount))
	router.HandleFunc("/accounts/{id}", makeHandleFunc(s.handleAccountById))
	router.HandleFunc("/transfer", makeHandleFunc(s.handleTrasfer))

	log.Println("JSON API Server running on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := getID(idStr)
	if err != nil {
		return err
	}
	if r.Method == "GET" {
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		id, err = s.store.DeleteAccount(id)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
		if id == 0 {
			err = fmt.Errorf("account %d not found", id)
			return WriteJSON(w, http.StatusNotFound, ApiError{Error: err.Error()})
		}
		return WriteJSON(w, http.StatusNoContent, map[string]int{"deleted": id})
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountRequest := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}

	account := NewAccount(createAccountRequest.FirstName, createAccountRequest.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, createAccountRequest)
}

func (s *ApiServer) handleTrasfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		transferRequest := &TransferRequest{}
		if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
			return err
		}
		id, err := s.store.Transfer(transferRequest.AccountNumber, transferRequest.Amount)
		if err != nil {
			return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
		if id == 0 {
			err = fmt.Errorf("account not %s not found", transferRequest.AccountNumber)
			return WriteJSON(w, http.StatusNotFound, ApiError{Error: err.Error()})
		}
		return WriteJSON(w, http.StatusOK, map[string]any{
			"transfered": transferRequest.Amount,
			"to":         transferRequest.AccountNumber,
		})
	}
	return fmt.Errorf("method %s not allowed", r.Method)
}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle errors in handle funcs
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func getID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
