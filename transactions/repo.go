package transactions

import (
	"database/sql"
	"errors"
	L "userTransactions/logging"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) getBalance(id int) float64 {
	var balance float64
	q := "SELECT balance FROM userbalance WHERE user_id = $1"
	err := r.db.QueryRow(q, id).Scan(&balance)
	if err != nil {
		L.Logger.Error("Failed getting user ", id, " balance", err)
		return 0
	}

	return balance
}

func (r *Repo) deposit(userId *int, amount float64) (float64, error) {
	if userId == nil {
		return 0, errors.New("user ID is required")
	}

	if amount <= 0 {
		return 0, errors.New("amount must be positive number")
	}

	var currentBalance float64
	q := "SELECT balance FROM userbalance WHERE user_id = $1"
	err := r.db.QueryRow(q, userId).Scan(&currentBalance)
	if err != nil {
		return 0, err
	}

	newBalance := currentBalance + amount
	q = "UPDATE userbalance SET balance = $1 WHERE user_id = $2"
	_, err = r.db.Exec(q, newBalance, userId)
	if err != nil {
		return 0, err
	}

	err = r.newTransactionRecord(userId, userId, "D", amount)

	return newBalance, err
}

func (r *Repo) transfer(senderId, recipientId *int, amount float64) error {
	if senderId == nil {
		return errors.New("sender ID is required")
	}

	if amount <= 0 {
		return errors.New("amount must be positive number")
	}

	if recipientId == nil {
		return errors.New("recipient is required to execute transaction")
	}

	var currentBalance float64
	q := "SELECT balance FROM userbalance WHERE user_id = $1"
	err := r.db.QueryRow(q, senderId).Scan(&currentBalance)
	if err != nil {
		return err
	}

	if currentBalance-amount >= 0 {
		return errors.New("not enough funds on the account")
	}

	newBalance := currentBalance - amount
	q = "UPDATE userbalance SET balance = $1 WHERE user_id = $2"
	_, err = r.db.Exec(q, newBalance, senderId)
	if err != nil {
		return err
	}

	var currentBalanceRecipient float64
	err = r.db.QueryRow(q, recipientId).Scan(&currentBalanceRecipient)
	if err != nil {
		r.restoreBalance(senderId, currentBalance)
		return err
	}

	newBalance = currentBalanceRecipient + amount
	q = "UPDATE userbalance SET balance = $1 WHERE user_id = $2"
	_, err = r.db.Exec(q, newBalance, senderId)
	if err != nil {
		r.restoreBalance(senderId, currentBalance)
		return err
	}

	err = r.newTransactionRecord(senderId, recipientId, "T", amount)
	if err != nil {
		r.restoreBalance(senderId, currentBalance)
		r.restoreBalance(recipientId, currentBalanceRecipient)
		return err
	}

	return nil
}

func (r *Repo) newTransactionRecord(from, to *int, operation string, amount float64) error {
	q := "INSERT INTO transaction (from, to, type, amount) VALUES ($1, $2, $3, $4)"

	_, err := r.db.Exec(q, from, to, operation, amount)
	return err
}

func (r *Repo) restoreBalance(userID *int, balance float64) {
	q := "UPDATE userbalance SET balance = $1 WHERE user_id = $2"
	r.db.Exec(q, balance, userID)
}
