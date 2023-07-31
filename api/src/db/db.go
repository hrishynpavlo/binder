package db

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func ProvideDb() *sqlx.DB {
	connStr := "user=binder_usr password=binder_best_app dbname=binder_all sslmode=disable"
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

type User struct {
	Id          int64          `db:"id"`
	Email       string         `db:"email"`
	FirstName   string         `db:"first_name"`
	LastName    string         `db:"last_name"`
	DisplayName string         `db:"display_name"`
	DateOfBirth string         `db:"date_of_birth"`
	Country     string         `db:"country"`
	Geolocation string         `db:"geolocation"`
	Interests   sql.NullString `db:"interests"`
}
