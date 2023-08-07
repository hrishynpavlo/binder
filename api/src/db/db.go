package db

import (
	"binder_api/configuration"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DB struct {
	*sqlx.DB
}

func ProvideDb(logger *zap.Logger, config *configuration.AppConfiguration) *sqlx.DB {
	db, err := sqlx.Connect("postgres", config.DbConnectionString)
	if err != nil {
		logger.Fatal("Server can't connect to dabase",
			zap.String("connection_string", config.DbConnectionString))
		panic("EXIT")
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

type UserMatchingInfo struct {
	UserId        int64         `db:"id"`
	Latitude      float64       `db:"latitude"`
	Longitude     float64       `db:"longitude"`
	Age           int8          `db:"age"`
	Interests     string        `db:"interests"`
	MaxDistanceKm sql.NullInt16 `db:"max_distance_km"`
	MinAge        sql.NullInt16 `db:"min_age"`
	MaxAge        sql.NullInt16 `db:"max_age"`
}

type UserWithActualGeo struct {
	Id                int64          `db:"id"`
	Email             string         `db:"email"`
	FirstName         string         `db:"first_name"`
	LastName          string         `db:"last_name"`
	DisplayName       string         `db:"display_name"`
	DateOfBirth       string         `db:"date_of_birth"`
	CountryCode       string         `db:"country_code"`
	StateCode         string         `db:"state_code"`
	City              string         `db:"city"`
	Latitude          float64        `db:"latitude"`
	Longitude         float64        `db:"longitude"`
	Interests         sql.NullString `db:"interests"`
	PhotoUrls         sql.NullString `db:"photo_urls"`
	PrimaryPhotoIndex sql.NullInt16  `db:"primary_photo_index"`
	MinDistanceKm     sql.NullInt16  `db:"min_distance_km"`
	MaxDistanceKm     sql.NullInt16  `db:"max_distance_km"`
	MinAge            sql.NullInt16  `db:"min_age"`
	MaxAge            sql.NullInt16  `db:"max_age"`
}
