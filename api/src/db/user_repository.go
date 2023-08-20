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

func (repo UserRepository) GetUserIdByEmailAndPassword(email string, passwordHash string) (int64, error) {
	result := UserId{}
	err := repo.db.Get(&result, "SELECT * FROM sp_login_user($1, $2)", email, passwordHash)
	if err != nil {
		repo.logger.Error("GetUserIdByEmailAndPassword() db error", zap.Error(err))
		return 0, err
	}

	return result.ID, nil
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

type UserId struct {
	ID int64 `db:"id"`
}
