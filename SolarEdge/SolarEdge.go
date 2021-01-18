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

func GetLatestPowerData(siteId int) (*sd.Power, error) {
	now := time.Now()
	then := now.Add(time.Duration(-15) * time.Minute)
	return getPowerData(siteId, now, then)
}

func GetPowerDataAt(siteId int, targetTime time.Time) (*sd.PowerReading, error) {
	start := targetTime.Add(-time.Duration(-5) * time.Minute)
	end := targetTime.Add(time.Duration(5) * time.Minute)
	pd, err := getPowerData(siteId, end, start)
	if err != nil {
		fmt.Sprintf("Error getting power data at [%s]: %s", targetTime, err)
		return nil, err
	}

	var retVal *sd.PowerReading

	for i := range pd.PowerDetails.Meters {
		if pd.PowerDetails.Meters[i].Type == "Production" {
			valueIndex := len(pd.PowerDetails.Meters[i].Values) - 1
			var pr *sd.PowerReading
			pr = new(sd.PowerReading)
			pr.Value = pd.PowerDetails.Meters[i].Values[valueIndex].Value
			layout := "2006-01-02 15:04:05"
			str := pd.PowerDetails.Meters[i].Values[valueIndex].Date
			t, timeParseErr := time.Parse(layout, str)

			if timeParseErr != nil {
				fmt.Printf("Error parsing power datetime: %s", timeParseErr)
			}
			loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
			pr.Date = t.In(loc).Add(5 * time.Hour)
			retVal = pr
			break
		}
	}

	return retVal, nil
}

func getPowerData(siteId int, end time.Time, start time.Time) (*sd.Power, error) {
	location, loadLocErr := time.LoadLocation("America/Indiana/Indianapolis")
	if loadLocErr != nil {
		log.Fatal(fmt.Sprintf("Error decoding loading time location response: %s", loadLocErr))
		return nil, loadLocErr
	}
	startLoc := start.In(location)
	endLoc := end.In(location)

	solarEdgeKey := os.Getenv("SOLAREDGE_APIKEY")

	startTime := fmt.Sprintf("%d-%02d-%d%%20%02d:%02d:%02d", startLoc.Year(), startLoc.Month(), startLoc.Day(), startLoc.Hour(), startLoc.Minute(), startLoc.Second())
	endTime := fmt.Sprintf("%d-%02d-%d%%20%02d:%02d:%02d", endLoc.Year(), endLoc.Month(), endLoc.Day(), endLoc.Hour(), endLoc.Minute(), endLoc.Second())

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
