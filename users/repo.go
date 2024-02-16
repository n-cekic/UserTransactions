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

func (r *Repo) getUserIDFromEmail(email string) (int, error) {
	q := `SELECT user_id FROM "user" WHERE email = $1`

	var id int
	err := r.db.QueryRow(q, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			L.Logger.Print("No matching record found for the given name.")
			return id, err
		}
	}
	return id, err
}
