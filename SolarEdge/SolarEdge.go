package SolarEdge

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sd "github.com/jasondborneman/solar3/SolarData"
)

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

func GetSolarSiteInfo(siteId int) (*sd.Site, error) {
	solarEdgeKey := os.Getenv("SOLAREDGE_APIKEY")
	url := fmt.Sprintf("https://monitoringapi.solaredge.com/site/%d/details?api_key=%s", siteId, solarEdgeKey)

	resp, callErr := callSolarWinds(url)
	if callErr != nil {
		log.Fatal(fmt.Sprintf("Error loading response for site data: %s", callErr))
		return nil, callErr
	}
	siteInfo := &sd.Site{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&siteInfo)
	if decodeErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding Site Info response: %s", decodeErr))
		return nil, decodeErr
	}
	return siteInfo, nil
}

func GetPastDayPowerData(siteId int) (*sd.Power, error) {
	now := time.Now()
	then := now.AddDate(0, 0, -1)
	return getPowerData(siteId, now, then)
}

func GetLatestPowerData(siteId int) (*sd.Power, error) {
	now := time.Now()
	then := now.Add(time.Duration(-15) * time.Minute)
	return getPowerData(siteId, now, then)
}

func getPowerData(siteId int, now time.Time, then time.Time) (*sd.Power, error) {
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
	powerInfo := &sd.Power{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&powerInfo)
	if decodeErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding Power Info response: %s", decodeErr))
		return nil, decodeErr
	}
	return powerInfo, nil
}
