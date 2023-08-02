package workers

import (
	"binder_api/configuration"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type GeoMatcherWorker struct {
	logger                *zap.Logger
	config                *configuration.AppConfiguration
	userRegisteredChannel chan UserRegisteredEvent
}

func ProvideGeoMatcherWorker(logger *zap.Logger, config *configuration.AppConfiguration, userRegisteredChannel chan UserRegisteredEvent) *GeoMatcherWorker {
	return &GeoMatcherWorker{logger: logger, config: config, userRegisteredChannel: userRegisteredChannel}
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
