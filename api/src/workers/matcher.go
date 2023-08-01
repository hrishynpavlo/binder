package workers

import (
	"binder_api/db"
	"context"
	"database/sql"
	"math"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	earthRadius = 6371 // Earth radius in km
)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func getOrDefault(value sql.NullInt16, defaultValue int16) int16 {
	if value.Valid {
		return value.Int16
	} else {
		return defaultValue
	}
}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	lat1Rad := degreesToRadians(lat1)
	lon1Rad := degreesToRadians(lon1)
	lat2Rad := degreesToRadians(lat2)
	lon2Rad := degreesToRadians(lon2)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

type MatcherWorker struct {
	logger *zap.Logger
	db     *sqlx.DB
	redis  *redis.Client
}

func ProvideMatcherWorker(logger *zap.Logger, db *sqlx.DB, redis *redis.Client) *MatcherWorker {
	return &MatcherWorker{logger: logger, db: db, redis: redis}
}

func (worker MatcherWorker) Start(ctx context.Context) {
	worker.logger.Info("Matcher worker started")
	userDistanceMap := make(map[int64][]int64)
	users := []db.UserMatchingInfo{}
	if err := worker.db.Select(&users, "select * from user_matching"); err != nil {
		worker.logger.Error("Getting users failed", zap.Error(err))
	}
	numberOfUsers := len(users)
	for index, user := range users {
		maxDistance := getOrDefault(user.MaxDistanceKm, 30000)
		for i := index + 1; i < numberOfUsers; i++ {
			possibleMatchUser := users[i]
			maxDistanceForPossibleUsers := getOrDefault(possibleMatchUser.MaxDistanceKm, 30000)
			distnace := getDistance(user.Latitude, user.Longitude, possibleMatchUser.Latitude, possibleMatchUser.Longitude)
			if maxDistance >= int16(distnace) {
				array := userDistanceMap[user.UserId]
				userDistanceMap[user.UserId] = append(array, possibleMatchUser.UserId)
			}

			if maxDistanceForPossibleUsers >= int16(distnace) {
				array := userDistanceMap[possibleMatchUser.UserId]
				userDistanceMap[possibleMatchUser.UserId] = append(array, user.UserId)
			}
		}
	}

	for key, value := range userDistanceMap {
		stringKey := strconv.FormatInt(key, 10)
		stringValue := make([]string, len(value))
		for i, id := range value {
			stringValue[i] = strconv.FormatInt(id, 10)
		}
		_, err := worker.redis.SAdd(context.Background(), stringKey, stringValue).Result()
		if err != nil {
			worker.logger.Error("", zap.Error(err))
		}
	}

	worker.logger.Info("", zap.Any("map", userDistanceMap))
}
