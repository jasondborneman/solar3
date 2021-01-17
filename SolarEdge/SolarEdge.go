package SolarEdge

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

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

func callSolarWinds(url string) (*http.Response, error) {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	resp, getErr := netClient.Get(url)
	if getErr != nil {
		log.Fatal(fmt.Sprintf("Error getting calling SolarEdge API: %s", getErr))
		return nil, getErr

	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Non-200 Status Code Returned: %d [%s]", resp.StatusCode, url)
		log.Fatal(message)
		return nil, errors.New(message)
	}
	return resp, nil
}

func GetSolarSiteInfo(siteId int) (*Site, error) {
	solarEdgeKey := os.Getenv("SOLAREDGE_APIKEY")
	url := fmt.Sprintf("https://monitoringapi.solaredge.com/site/%d/details?api_key=%s", siteId, solarEdgeKey)

	resp, callErr := callSolarWinds(url)
	if callErr != nil {
		log.Fatal(fmt.Sprintf("Error loading response for site data: %s", callErr))
		return nil, callErr
	}
	siteInfo := &Site{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&siteInfo)
	if decodeErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding Site Info response: %s", decodeErr))
		return nil, decodeErr
	}
	return siteInfo, nil
}

func GetPowerData(siteId int) (*Power, error) {
	now := time.Now()
	then := now.Add(time.Duration(-15) * time.Minute)

	location, loadLocErr := time.LoadLocation("America/Indiana/Indianapolis")
	if loadLocErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding loading time location response: %s", loadLocErr))
		return nil, loadLocErr
	}
	thenLoc := then.In(location)
	nowLoc := now.In(location)

	solarEdgeKey := os.Getenv("SOLAREDGE_APIKEY")

	startTime := fmt.Sprintf("%d-%02d-%d%%20%02d:%02d:%02d", thenLoc.Year(), thenLoc.Month(), thenLoc.Day(), thenLoc.Hour(), thenLoc.Minute(), thenLoc.Second())
	endTime := fmt.Sprintf("%d-%02d-%d%%20%02d:%02d:%02d", nowLoc.Year(), nowLoc.Month(), nowLoc.Day(), nowLoc.Hour(), nowLoc.Minute(), nowLoc.Second())

	url := fmt.Sprintf("https://monitoringapi.solaredge.com/site/%d/powerDetails.json?startTime=%s&endTime=%s&api_key=%s", siteId, startTime, endTime, solarEdgeKey)

	resp, callErr := callSolarWinds(url)
	if callErr != nil {
		log.Fatal(fmt.Sprintf("Error loading response for power data: %s", callErr))
		return nil, callErr
	}
	powerInfo := &Power{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&powerInfo)
	if decodeErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding Power Info response: %s", decodeErr))
		return nil, decodeErr
	}
	return powerInfo, nil
}
