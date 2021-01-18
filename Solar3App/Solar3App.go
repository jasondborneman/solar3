package Solar3App

import (
	"fmt"
	"os"
	"strconv"
	"time"

	g "github.com/jasondborneman/solar3/Graphing"
	ipgl "github.com/jasondborneman/solar3/IPGeoLocation"
	ow "github.com/jasondborneman/solar3/OpenWeather"
	s3data "github.com/jasondborneman/solar3/Solar3DataStorage"
	sd "github.com/jasondborneman/solar3/SolarData"
	se "github.com/jasondborneman/solar3/SolarEdge"
	tw "github.com/jasondborneman/solar3/Twitter"
)

func GetData() sd.SolarData {
	var retVal sd.SolarData

	now := time.Now()
	var siteID int
	siteID, _ = strconv.Atoi(os.Getenv("SOLAREDGE_SITEID"))
	latitude := os.Getenv("SOLAR3_LATITUDE")
	longitude := os.Getenv("SOLAR3_LONGITUDE")
	var site *sd.Site
	var powerNow *sd.Power
	var sunMoon *ipgl.SunMoonInfo
	var openWeather *ow.OpenWeather
	site, _ = se.GetSolarSiteInfo(siteID)
	powerNow, _ = se.GetLatestPowerData(siteID)
	sunMoon, _ = ipgl.GetSunPosition(latitude, longitude)
	openWeather, _ = ow.GetWeather(latitude, longitude)
	retVal.Site = site
	retVal.DateTimeStored = now
	for i := range powerNow.PowerDetails.Meters {
		if powerNow.PowerDetails.Meters[i].Type == "Production" {
			valueIndex := len(powerNow.PowerDetails.Meters[i].Values) - 1
			var pr sd.PowerReading
			pr.Value = powerNow.PowerDetails.Meters[i].Values[valueIndex].Value
			layout := "2006-01-02 15:04:05"
			str := powerNow.PowerDetails.Meters[i].Values[valueIndex].Date
			t, timeParseErr := time.Parse(layout, str)

			if timeParseErr != nil {
				fmt.Printf("Error parsing power datetime: %s", timeParseErr)
			}
			loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
			pr.Date = t.In(loc).Add(5 * time.Hour)
			retVal.PowerGen = pr
			retVal.DateTime = pr.Date
			break
		}
	}
	retVal.CloudCover = openWeather.Clouds.All
	retVal.Temp = openWeather.Main.Temp
	retVal.Pressure = openWeather.Main.Pressure
	retVal.Humidity = openWeather.Main.Humidity
	retVal.SunAzimuth = sunMoon.SunAzimuth
	retVal.SunAltitude = sunMoon.SunAltitude
	retVal.SnowOneHr = openWeather.Snow.OneH
	retVal.RainOneHr = openWeather.Rain.OneH
	retVal.WeatherID = openWeather.Weather[0].ID
	return retVal
}

func Run(doTweet bool, doSaveGraph bool) {
	var data sd.SolarData
	data = GetData()
	xVals, powerYVals, cloudYVals, maxPower, errSave := s3data.SaveToFirestore(data)
	saved := true
	if errSave != nil {
		saved = false
	}
	graphBytes := g.CreateGraph(xVals, powerYVals, cloudYVals, maxPower)
	pngSaved := false
	if doSaveGraph {
		errSaveJpeg := g.SaveGraph(graphBytes, "chart")
		pngSaved = true
		if errSaveJpeg != nil {
			pngSaved = false
		}
	}
	message := `
DateTime: %s
Last Reported Power: %.2f
Cloud Cover: %d
Sun Azimuth: %.2f
Sun Altitude: %.2f`
	loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
	generatedDate := data.PowerGen.Date.In(loc)
	message = fmt.Sprintf(
		message,
		fmt.Sprintf("%02d-%02d-%d %02d:%02d",
			generatedDate.Month(),
			generatedDate.Day(),
			generatedDate.Year(),
			generatedDate.Hour(),
			generatedDate.Minute()),
		data.PowerGen.Value,
		data.CloudCover,
		data.SunAzimuth,
		data.SunAltitude)
	tweeted := false
	if doTweet {
		tweetErr := tw.TweetWithMedia(message, graphBytes)
		tweeted = true
		if tweetErr != nil {
			tweeted = false
		}
	}
	dodgyTimes := s3data.GetDodgyDataTimesPast24Hrs()
	siteID, _ := strconv.Atoi(os.Getenv("SOLAREDGE_SITEID"))
	fixedCount := 0
	fixErrCount := 0
	for i := range dodgyTimes {
		t := dodgyTimes[i]
		fixedPower, gpdaErr := se.GetPowerDataAt(siteID, t)
		if gpdaErr != nil {
			fmt.Printf("Error getting correct power data at [%s]: %s", t.In(loc), gpdaErr)
			fixErrCount++
		}
		updateErr := s3data.UpdatePowerDataAt(fixedPower)
		if updateErr != nil {
			fmt.Printf("Error updating dodgy power data at [%s]: %s", t.In(loc), updateErr)
		}
		fixedCount++
	}
	fmt.Println("----------------------------")
	fmt.Printf("Saved To Firestore?: %t\n", saved)
	fmt.Printf("Dodgy Data Count:    %d\n", len(dodgyTimes))
	fmt.Printf("Fixed Data Count:    %d\n", fixedCount)
	fmt.Printf("Fix Data Err Count:  %d\n", fixErrCount)
	fmt.Printf("Saved Graph?:        %t\n", pngSaved)
	fmt.Printf("Tweeted?:            %t\n", tweeted)
	fmt.Println(message)
}
