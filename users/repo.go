package users

import (
	"database/sql"
	L "userTransactions/logging"

	_ "github.com/lib/pq"
)

type Repo struct {
	db *sql.DB
}

func (r *Repo) createUser(email string) (int, error) {
	q := `INSERT INTO "user" ("email") VALUES ($1) RETURNING user_id`

	var id int
	err := r.db.QueryRow(q, email).Scan(&id)
	if err != nil {
		L.Logger.Errorf("New user insert failed: %s", err.Error())
		return 0, err
	}
	return id, nil
}

func (r *Repo) getUserIDFromEmail(email string) (int, error) {
	q := `SELECT user_id FROM "user" WHERE email = $1`

	var id int
	err := r.db.QueryRow(q, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			L.Logger.Info("No matching record found for the given email: ", email)
			return id, err
		}
	}
	return id, err
}
