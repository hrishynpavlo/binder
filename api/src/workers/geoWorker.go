package workers

import (
	"binder_api/configuration"
	"binder_api/db"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type GeoWorker struct {
	logger                *zap.Logger
	config                *configuration.AppConfiguration
	repo                  *db.UserRepository
	userRegisteredChannel chan UserRegisteredEvent
	redis                 *redis.Client
}

func ProvideGeoMatcherWorker(logger *zap.Logger, config *configuration.AppConfiguration, userRegisteredChannel chan UserRegisteredEvent, redis *redis.Client, repo *db.UserRepository) *GeoWorker {
	return &GeoWorker{logger: logger, config: config, userRegisteredChannel: userRegisteredChannel, redis: redis, repo: repo}
}

func ProvideUserRegisteredChannel() chan UserRegisteredEvent {
	return make(chan UserRegisteredEvent)
}

type UserRegisteredEvent struct {
	UserId    int64
	Latitude  float64
	Longitude float64
}

func (worker GeoWorker) StartGeoEnrichment() {
	for message := range worker.userRegisteredChannel {
		worker.logger.Info("Geo worker process new event", zap.Any("userRegisteredEvent", message))

		location, err := worker.getReverseGeocode(message.Latitude, message.Longitude)

		if err != nil {
			worker.logger.Error("Geo worker HERE request error", zap.Error(err))
			continue
		}

		worker.logger.Info("HERE api response", zap.Any("location", location))

		if len(location.Items) < 1 {
			continue
		}

		if err := worker.repo.UpdateUserGeo(message.UserId, location.Items[0].Address.CountryCode, location.Items[0].Address.StateCode, location.Items[0].Address.City, message.Latitude, message.Longitude); err != nil {
			worker.logger.Error("Geo worker error during writing in db", zap.Error(err))
		}

		countryKey := location.Items[0].Address.CountryCode
		stateKey := location.Items[0].Address.StateCode
		city := strings.ReplaceAll(location.Items[0].Address.City, " ", "")
		cityKey := fmt.Sprintf("%s:%s:%s", countryKey, stateKey, city)
		userId := strconv.FormatInt(message.UserId, 10)

		worker.redis.SAdd(context.Background(), countryKey, userId).Result()
		worker.redis.SAdd(context.Background(), cityKey, userId).Result()
	}
}

type ReverseGeocodeResponse struct {
	Items []ReverseGeocodeResponseItem `json:"items"`
}

type ReverseGeocodeResponseItem struct {
	Title   string                            `json:"title"`
	Id      string                            `json:"id"`
	Address ReverseGeocodeResponseItemAddress `json:"address"`
}

type ReverseGeocodeResponseItemAddress struct {
	Label       string `json:"label"`
	CountryCode string `json:"countryCode"`
	CountryName string `json:"countryName"`
	State       string `json:"state"`
	StateCode   string `json:"stateCode"`
	City        string `json:"city"`
}

func (worker GeoWorker) getReverseGeocode(latitude float64, longitude float64) (ReverseGeocodeResponse, error) {
	url := fmt.Sprintf("https://revgeocode.search.hereapi.com/v1/revgeocode?at=%f,%f&lang=en-US&apiKey=%s", latitude, longitude, worker.config.BinderHereGeoToken)
	statusCode, body, err := fasthttp.Get(nil, url)
	if err != nil {
		worker.logger.Error("Geo worker request to HERE api failed", zap.Error(err))
		return ReverseGeocodeResponse{}, err
	}

	if statusCode != 200 {
		worker.logger.Error("Geo worker request not 200", zap.Any("error", body))
		return ReverseGeocodeResponse{}, errors.New(string(body))
	}
	var location ReverseGeocodeResponse
	if err := json.Unmarshal(body, &location); err != nil {
		worker.logger.Error("Geo worker error during parsing response json", zap.Error(err), zap.Any("body", string(body)))
		return ReverseGeocodeResponse{}, err
	}

	return location, nil
}
