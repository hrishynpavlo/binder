package workers

import (
	"binder_api/configuration"
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

type GeoMatcherWorker struct {
	logger                *zap.Logger
	config                *configuration.AppConfiguration
	userRegisteredChannel chan UserRegisteredEvent
	redis                 *redis.Client
}

func ProvideGeoMatcherWorker(logger *zap.Logger, config *configuration.AppConfiguration, userRegisteredChannel chan UserRegisteredEvent, redis *redis.Client) *GeoMatcherWorker {
	return &GeoMatcherWorker{logger: logger, config: config, userRegisteredChannel: userRegisteredChannel, redis: redis}
}

func ProvideUserRegisteredChannel() chan UserRegisteredEvent {
	return make(chan UserRegisteredEvent)
}

type UserRegisteredEvent struct {
	UserId    int64
	Latitude  float64
	Longitude float64
}

func (worker GeoMatcherWorker) StartWorker() {
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

		countryKey := location.Items[0].Address.CountryCode
		city := strings.ReplaceAll(location.Items[0].Address.City, " ", "")
		cityKey := fmt.Sprintf("%s:%s", countryKey, city)
		userId := strconv.FormatInt(message.UserId, 10)

		worker.redis.SAdd(context.Background(), countryKey, userId)
		worker.redis.SAdd(context.Background(), cityKey, userId)
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
	City        string `json:"city"`
}

func (worker GeoMatcherWorker) getReverseGeocode(latitude float64, longitude float64) (ReverseGeocodeResponse, error) {
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
