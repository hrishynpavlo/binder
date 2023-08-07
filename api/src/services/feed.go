package services

import (
	"binder_api/db"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
)

type FeedProvider interface {
	GetOrAdd(userId int64) []FeedEntry
}

func ProvideFeedService(redis *redis.Client, repo *db.FeedRepository) FeedProvider {
	return &FeedRedisDbDecoratedProvider{redis: redis, repo: repo}
}

type FeedRedisDbDecoratedProvider struct {
	redis *redis.Client
	repo  *db.FeedRepository
}

func (feed FeedRedisDbDecoratedProvider) GetOrAdd(userId int64) []FeedEntry {
	var userFeed []FeedEntry
	key := fmt.Sprintf("feed:%d", userId)
	value, _ := feed.redis.Get(context.Background(), key).Result()

	// if no feed in cache then find in db
	if value == "" {
		snapshot, _ := feed.repo.GetUserFeedSnapshot(userId)

		// if no feed snapshot in db then calculate it
		if snapshot.Feed == "" {
			users, _ := feed.repo.FindUserFeed(userId)
			userMap := make(map[int64]db.UserWithActualGeoDTO)
			for _, user := range users {
				userMap[user.UserId] = user
			}
			distnaceMap := calculateFeed(userId, users)
			for key, distance := range distnaceMap {
				feedEntry := FeedEntry{User: userMap[key], DistanceInKm: distance}
				userFeed = append(userFeed, feedEntry)
			}

			bytes, _ := json.Marshal(userFeed)
			value = string(bytes)

			currentUser := userMap[userId]
			feed.repo.CreateFeedSnapshot(userId, currentUser.CountryCode, currentUser.StateCode, currentUser.City, value)
		} else {
			value = snapshot.Feed
			json.Unmarshal([]byte(value), &userFeed)
		}

		feed.redis.Set(context.Background(), key, value, 1*time.Hour).Result()
	} else {
		json.Unmarshal([]byte(value), &userFeed)
	}

	return userFeed
}

func calculateFeed(userId int64, users []db.UserWithActualGeoDTO) map[int64]float64 {
	distanceMap := make(map[int64]float64)
	numberOfUsers := len(users)
	for index, user := range users {
		for i := index + 1; i < numberOfUsers; i++ {
			userToCalculate := users[i]
			if user.UserId != userId && userToCalculate.UserId != userId {
				continue
			}

			var mapId int64
			if user.UserId == userId {
				mapId = userToCalculate.UserId
			} else {
				mapId = user.UserId
			}

			distnace := getDistance(user.Latitude, user.Longitude, userToCalculate.Latitude, userToCalculate.Longitude)

			distanceMap[mapId] = distnace
		}
	}
	return distanceMap
}

const (
	earthRadius = 6371 // Earth radius in km
)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
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

func mapFeedEntry(user db.UserWithActualGeoDTO, distance float64) FeedEntry {
	return FeedEntry{
		User:         user,
		DistanceInKm: distance,
	}
}

type FeedEntry struct {
	User         db.UserWithActualGeoDTO
	DistanceInKm float64
}
