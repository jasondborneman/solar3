package IPGeoLocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type SunMoonInfo struct {
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Date                 string  `json:"date"`
	CurrentTime          string  `json:"current_time"`
	Sunrise              string  `json:"sunrise"`
	Sunset               string  `json:"sunset"`
	SunStatus            string  `json:"sun_status"`
	SolarNoon            string  `json:"solar_noon"`
	DayLength            string  `json:"day_length"`
	SunAltitude          float64 `json:"sun_altitude"`
	SunDistance          float64 `json:"sun_distance"`
	SunAzimuth           float64 `json:"sun_azimuth"`
	Moonrise             string  `json:"moonrise"`
	Moonset              string  `json:"moonset"`
	MoonStatus           string  `json:"moon_status"`
	MoonAltitude         float64 `json:"moon_altitude"`
	MoonDistance         float64 `json:"moon_distance"`
	MoonAzimuth          float64 `json:"moon_azimuth"`
	MoonParallacticAngle float64 `json:"moon_parallactic_angle"`
}

func callIPGeoLocation(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, getErr := netClient.Get(url)
	if getErr != nil {
		log.Fatalf("Error getting calling IPGeoLocation API: %s", getErr)
		return nil, getErr

	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Non-200 Status Code Returned: %d [%s]", resp.StatusCode, url)
		log.Fatal(message)
		return nil, errors.New(message)
	}
	return resp, nil
}

func GetSunPosition(latitude string, longitude string) (*SunMoonInfo, error) {
	ipGeoLocationKey := os.Getenv("IPGEOLOCATION_APIKEY")
	url := fmt.Sprintf("https://api.ipgeolocation.io/astronomy?apiKey=%s&lat=%s&long=%s", ipGeoLocationKey, latitude, longitude)

	resp, callErr := callIPGeoLocation(url)
	if callErr != nil {
		log.Fatalf("Error loading response for sun position: %s", callErr)
		return nil, callErr
	}
	sunMoonInfo := &SunMoonInfo{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&sunMoonInfo)
	if decodeErr != nil {
		log.Fatalf("Error decoding sun position response: %s", decodeErr)
		return nil, decodeErr
	}
	return sunMoonInfo, nil
}
