package AirNow

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type AQI []struct {
	DateObserved  string  `json:"DateObserved"`
	HourObserved  int     `json:"HourObserved"`
	LocalTimeZone string  `json:"LocalTimeZone"`
	ReportingArea string  `json:"ReportingArea"`
	StateCode     string  `json:"StateCode"`
	Latitude      float64 `json:"Latitude"`
	Longitude     float64 `json:"Longitude"`
	ParameterName string  `json:"ParameterName"`
	Aqi           int     `json:"AQI"`
	Category      struct {
		Number int    `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Category"`
}

type AQIData struct {
	PM2_5AQI int
	PM10AQI  int
}

func callAirNow(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, getErr := netClient.Get(url)
	if getErr != nil {
		log.Fatalf("Error getting calling AirNow API: %s", getErr)
		return nil, getErr

	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Non-200 Status Code Returned: %d [%s]", resp.StatusCode, url)
		log.Fatal(message)
		return nil, errors.New(message)
	}
	return resp, nil
}

func GetAQI(latitude, longitude string) (*AQIData, error) {
	airNowAPIKey := os.Getenv("AIRNOW_APIKEY")
	url := fmt.Sprintf("http://www.airnowapi.org/aq/observation/latLong/current/?format=application/json&latitude=%s&longitude=%s&distance=25&API_KEY=%s", latitude, longitude, airNowAPIKey)

	resp, callErr := callAirNow(url)
	if callErr != nil {
		log.Fatalf("Error loading response for aqi: %s", callErr)
		return nil, callErr
	}
	aqiInfo := &AQI{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&aqiInfo)
	if decodeErr != nil {
		log.Fatalf("Error decoding aqi response: %s", decodeErr)
		return nil, decodeErr
	}
	pm2_5 := 0
	pm10 := 0
	for _, a := range *aqiInfo {
		if a.ParameterName == "PM2.5" {
			pm2_5 = a.Aqi
		} else if a.ParameterName == "PM10" {
			pm10 = a.Aqi
		}
	}
	retVal := &AQIData{
		PM2_5AQI: pm2_5,
		PM10AQI:  pm10,
	}
	return retVal, nil
}
