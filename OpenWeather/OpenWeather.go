package OpenWeather

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type OpenWeather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Snow struct {
		OneH float64 `json:"1h"`
	} `json:"snow"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func callOpenWeather(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, getErr := netClient.Get(url)
	if getErr != nil {
		log.Fatal(fmt.Sprintf("Error getting calling OpenWeather API: %s", getErr))
		return nil, getErr

	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Non-200 Status Code Returned: %d [%s]", resp.StatusCode, url)
		log.Fatal(message)
		return nil, errors.New(message)
	}
	return resp, nil
}

func GetWeather(latitude string, longitude string) (*OpenWeather, error) {
	openWeatherAPIKey := os.Getenv("OPENWEATHER_APIKEY")
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=imperial&appid=%s", latitude, longitude, openWeatherAPIKey)

	resp, callErr := callOpenWeather(url)
	if callErr != nil {
		log.Fatal(fmt.Sprintf("Error loading response for sun position: %s", callErr))
		return nil, callErr
	}
	openWeatherInfo := &OpenWeather{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&openWeatherInfo)
	if decodeErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding sun position response: %s", decodeErr))
		return nil, decodeErr
	}
	return openWeatherInfo, nil
}
