package db

import (
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type UserRepository struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func ProvideUserRepository(logger *zap.Logger, db *sqlx.DB) *UserRepository {
	return &UserRepository{logger: logger, db: db}
}

func (repo UserRepository) GetAllUsers() ([]UserDTO, error) {
	var dbUsers []User
	err := repo.db.Select(&dbUsers, "SELECT * FROM users_info")
	if err != nil {
		repo.logger.Error("GetAllUsers() sql error", zap.Error(err))
		return []UserDTO{}, err
	}
	users := make([]UserDTO, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = mapUser(dbUser)
	}
	return users, nil
}

func (repo UserRepository) GetUserById(userId int64) (UserDTO, error) {
	var user User
	err := repo.db.Get(&user, "SELECT * FROM users_info where id = $1", userId)
	if err != nil {
		repo.logger.Error("GetUserById() error", zap.Error(err))
		return UserDTO{}, errors.New("database error")
	}

	return mapUser(user), nil
}

func (repo UserRepository) CreateUser(email string, passwordHash string, firstName string, lastName string, displayName string, dateOfBirth string, country string, latitude float64, longitude float64) (UserDTO, error) {
	user := User{}
	err := repo.db.Get(&user, "SELECT * FROM sp_create_user($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		email, passwordHash, firstName, lastName, displayName, dateOfBirth, country, latitude, longitude)
	if err != nil {
		repo.logger.Error("CreateUser() failed",
			zap.Error(err),
			zap.String("user_email", email))

		return UserDTO{}, errors.New("database error")
	}

	return mapUser(user), nil
}

func (repo UserRepository) UpdateUserInterests(userId int64, interests []string) (UserDTO, error) {
	var user User
	if err := repo.db.Get(&user, "SELECT * FROM sp_update_user_interests($1, $2)",
		userId, pq.Array(interests)); err != nil {
		repo.logger.Error("UpdateUserInterests() failed", zap.Error(err))
		return UserDTO{}, errors.New("database error")
	}

	return mapUser(user), nil
}

func (repo UserRepository) UpdateUserPhoto(userId int64, photoUrls []string) (UserDTO, error) {
	var user User
	if err := repo.db.Get(&user, "SELECT * FROM sp_update_user_photos($1, $2)", userId, pq.Array(photoUrls)); err != nil {

		repo.logger.Error("UpdateUserPhoto() failed", zap.Error(err))
		return UserDTO{}, errors.New("database error")
	}

	return mapUser(user), nil
}

func (repo UserRepository) UpdateUserFilter(userId int64, minDistanceKm int8, maxDistanceKm int8, minAge int8, maxAge int8) (UserDTO, error) {
	var user User
	if err := repo.db.Get(&user, "SELECT * FROM sp_update_user_filters($1, $2, $3, $4, $5)", userId, minDistanceKm, maxDistanceKm, minAge, maxAge); err != nil {
		repo.logger.Error("UpdateUserFilter() failed", zap.Error(err))
		return UserDTO{}, errors.New("database error")
	}

	return mapUser(user), nil
}

func (repo UserRepository) UpdateUserGeo(userId int64, countryCode string, stateCode string, city string, latitude float64, longitude float64) error {
	_, err := repo.db.Exec("select sp_update_user_geo($1, $2, $3, $4, $5, $6)", userId, countryCode, stateCode, city, latitude, longitude)
	if err != nil {
		repo.logger.Error("UpdateUserGeo() failed", zap.Error(err), zap.Int64("user_id", userId))
		return err
	}

	return nil
}

func (repo UserRepository) GetUserFeed(userId int64) ([]UserDTO, error) {
	var dbUsers []User
	if err := repo.db.Select(&dbUsers, "with user_location as (select ug.country_code, ug.state_code, ug.city from user_geos ug where ug.user_id = 52), feed_users as (select ug.user_id, ug.geolocation from user_geos ug where ug.country_code in (select ul.country_code from user_location ul) and ug.state_code in (select ul.state_code from user_location ul) and ug.city in (select ul.city from user_location ul) and ug.user_id <> 52) select ui.id, ui.email, ui.first_name, ui.last_name, ui.display_name, ui.date_of_birth, ui.country, fu.geolocation, ui.interests, ui.photo_urls, ui.primary_photo_index, ui.min_distance_km, ui.max_distance_km, ui.min_age, ui.max_age from users_info ui join feed_users fu on fu.user_id = ui.id;", userId); err != nil {
		return nil, err
	}

	users := make([]UserDTO, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = mapUser(dbUser)
	}
	return users, nil
}

type Interest string

const (
	TRAVEL    Interest = "Travel"
	MUSIC     Interest = "Music"
	BOOKS     Interest = "Books"
	MOVIES    Interest = "Movies"
	SPORT     Interest = "Sport"
	ADVENTURE Interest = "Adventure"
	PETS      Interest = "Pets"
	ANIMALS   Interest = "Animals"
	FOOD      Interest = "Food"
	WINE      Interest = "Wine"
	COFFEE    Interest = "Coffee"
	DRINK     Interest = "Drink"
	WALKS     Interest = "Walks"
	HIKING    Interest = "Hiking"
	DANCING   Interest = "Dancing"
	GYM       Interest = "Gym"
	TATTOO    Interest = "Tattoo"
)

type UserDTO struct {
	Id                int64
	Email             string
	FirstName         string
	LastName          string
	DisplayName       string
	DateOfBirth       string
	Interests         []Interest
	PhotoUrls         []string
	PrimaryPhotoIndex int16
	MaxDistanceKm     int16
	MinAge            int16
	MaxAge            int16
}

func mapUser(user User) UserDTO {
	var interests []Interest
	var photoUrls []string
	primaryPhotoIndex := int16(0)
	maxDistanceKm := int16(10)
	minAge := int16(18)
	maxAge := int16(25)
	if user.Interests.Valid {
		interests = parseInterests(user.Interests.String)
	}

	if user.PhotoUrls.Valid {
		photoUrls = parsePgArray(user.PhotoUrls.String)
	}

	if user.PrimaryPhotoIndex.Valid {
		primaryPhotoIndex = user.PrimaryPhotoIndex.Int16
	}

	if user.MaxDistanceKm.Valid {
		maxDistanceKm = user.MinDistanceKm.Int16
	}

	if user.MinAge.Valid {
		minAge = user.MinAge.Int16
	}

	if user.MaxAge.Valid {
		maxAge = user.MaxAge.Int16
	}

	return UserDTO{Id: user.Id, Email: user.Email, FirstName: user.FirstName, LastName: user.LastName,
		DisplayName: user.DisplayName, DateOfBirth: user.DateOfBirth, Interests: interests, PhotoUrls: photoUrls,
		PrimaryPhotoIndex: primaryPhotoIndex, MaxDistanceKm: maxDistanceKm, MinAge: minAge, MaxAge: maxAge}
}

func parsePgArray(original string) []string {
	original = strings.ReplaceAll(original, "{", "")
	original = strings.ReplaceAll(original, "}", "")
	result := strings.Split(original, ",")
	return result
}

func parseInterests(original string) []Interest {
	strs := parsePgArray(original)
	interests := make([]Interest, len(strs))
	for i, str := range strs {
		interests[i] = Interest(str)
	}
	return interests
}
