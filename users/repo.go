package users

import (
	"database/sql"
	L "userTransactions/logging"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) createUser(email string) error {
	q := `INSERT INTO "user" ("email") VALUES ($1)`

	_, err := r.db.Exec(q, email)
	if err != nil {
		L.Logger.Printf("New user insert failed: %s", err.Error())
	}
	return err
}
