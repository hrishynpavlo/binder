package db

import (
	"binder_api/configuration"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func ProvideDb(config *configuration.AppConfiguration) *sqlx.DB {
	db, err := sqlx.Connect("postgres", config.DbConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

type User struct {
	Id                int64          `db:"id"`
	Email             string         `db:"email"`
	FirstName         string         `db:"first_name"`
	LastName          string         `db:"last_name"`
	DisplayName       string         `db:"display_name"`
	DateOfBirth       string         `db:"date_of_birth"`
	Country           string         `db:"country"`
	Geolocation       string         `db:"geolocation"`
	Interests         sql.NullString `db:"interests"`
	PhotoUrls         sql.NullString `db:"photo_urls"`
	PrimaryPhotoIndex sql.NullInt16  `db:"primary_photo_index"`
	MinDistanceKm     sql.NullInt16  `db:"min_distance_km"`
	MaxDistanceKm     sql.NullInt16  `db:"max_distance_km"`
	MinAge            sql.NullInt16  `db:"min_age"`
	MaxAge            sql.NullInt16  `db:"max_age"`
}
