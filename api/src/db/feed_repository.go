package db

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type FeedRepository struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func ProvideFeedRepository(logger *zap.Logger, db *sqlx.DB) *FeedRepository {
	return &FeedRepository{logger: logger, db: db}
}

func (repo FeedRepository) GetUserFeedSnapshot(userId int64) (UserFeedSnapshot, error) {
	var snapshot UserFeedSnapshot
	if err := repo.db.Get(&snapshot, "SELECT f.feed, f.feed_offset FROM user_feeds f WHERE f.user_id = $1 AND is_active = TRUE LIMIT 1;", userId); err != nil {
		repo.logger.Error("GetUserFeedSnapshot() failed", zap.Error(err), zap.Int64("user_id", userId))
		return snapshot, err
	}

	return snapshot, nil
}

func (repo FeedRepository) FindUserFeed(userId int64) ([]UserWithActualGeoDTO, error) {
	var dbUsers []UserWithActualGeo
	if err := repo.db.Select(&dbUsers, "SELECT * FROM sp_find_users_to_match($1);", userId); err != nil {
		repo.logger.Error("GetUserFeed() failed", zap.Error(err), zap.Int64("user_id", userId))
		return nil, err
	}

	users := make([]UserWithActualGeoDTO, len(dbUsers))
	for index, dbUser := range dbUsers {
		users[index] = dbUser.toDto()
	}

	return users, nil
}

func (repo FeedRepository) CreateFeedSnapshot(userId int64, countryCode string, stateCode string, city string, feedJson string) error {
	_, err := repo.db.Exec("SELECT sp_create_feed_snapshot($1, $2, $3, $4, $5);", userId, countryCode, stateCode, city, feedJson)
	if err != nil {
		repo.logger.Error("CreateFeedSnapshot() failed", zap.Error(err), zap.Int64("user_id", userId))
	}
	return err
}

type UserFeedSnapshot struct {
	Feed   string `db:"feed"`
	Offset int8   `db:"feed_offset"`
}

type UserWithActualGeoDTO struct {
	UserId      int64
	FirstName   string
	LastName    string
	DisplayName string
	DateOfBirth string
	CountryCode string
	StateCode   string
	City        string
	Latitude    float64
	Longitude   float64
	Interests   []Interest
	PhotoUrls   []string
}

func (record UserWithActualGeo) toDto() UserWithActualGeoDTO {
	interests := parseInterests(record.Interests.String)
	photoUrsl := parsePgArray(record.PhotoUrls.String)
	return UserWithActualGeoDTO{UserId: record.Id, FirstName: record.FirstName, LastName: record.LastName, DisplayName: record.DisplayName,
		DateOfBirth: record.DateOfBirth, CountryCode: record.CountryCode, StateCode: record.StateCode, City: record.City,
		Latitude: record.Latitude, Longitude: record.Longitude, Interests: interests, PhotoUrls: photoUrsl,
	}
}
