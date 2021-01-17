package SolarData

import (
	"time"
)

type SolarData struct {
	Site           *Site
	DateTime       time.Time
	DateTimeStored time.Time
	PowerGen       PowerReading
	CloudCover     int
	Temp           float64
	Pressure       int
	Humidity       int
	SunAzimuth     float64
	SunAltitude    float64
	SnowOneHr      float64
	RainOneHr      float64
	WeatherID      int
}

type Site struct {
	Details struct {
		ID               int         `json:"id"`
		Name             string      `json:"name"`
		AccountID        int         `json:"accountId"`
		Status           string      `json:"status"`
		PeakPower        float64     `json:"peakPower"`
		LastUpdateTime   string      `json:"lastUpdateTime"`
		InstallationDate string      `json:"installationDate"`
		PtoDate          interface{} `json:"ptoDate"`
		Notes            string      `json:"notes"`
		Type             string      `json:"type"`
		Location         struct {
			Country     string `json:"country"`
			State       string `json:"state"`
			City        string `json:"city"`
			Address     string `json:"address"`
			Address2    string `json:"address2"`
			Zip         string `json:"zip"`
			TimeZone    string `json:"timeZone"`
			CountryCode string `json:"countryCode"`
			StateCode   string `json:"stateCode"`
		} `json:"location"`
		PrimaryModule struct {
			ManufacturerName string  `json:"manufacturerName"`
			ModelName        string  `json:"modelName"`
			MaximumPower     float64 `json:"maximumPower"`
		} `json:"primaryModule"`
		Uris struct {
			SITEIMAGE      string `json:"SITE_IMAGE"`
			DATAPERIOD     string `json:"DATA_PERIOD"`
			INSTALLERIMAGE string `json:"INSTALLER_IMAGE"`
			DETAILS        string `json:"DETAILS"`
			OVERVIEW       string `json:"OVERVIEW"`
		} `json:"uris"`
		PublicSettings struct {
			IsPublic bool `json:"isPublic"`
		} `json:"publicSettings"`
	} `json:"details"`
}

type Power struct {
	PowerDetails struct {
		TimeUnit string `json:"timeUnit"`
		Unit     string `json:"unit"`
		Meters   []struct {
			Type   string `json:"type"`
			Values []struct {
				Date  string  `json:"date"`
				Value float64 `json:"value"`
			} `json:"values"`
		} `json:"meters"`
	} `json:"powerDetails"`
}

type PowerReading struct {
	Date  time.Time
	Value float64
}
