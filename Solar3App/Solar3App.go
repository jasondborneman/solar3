package Solar3App

import (
	"fmt"
	ipgl "github.com/jasondborneman/solar3/IPGeoLocation"
	ow "github.com/jasondborneman/solar3/OpenWeather"
	s3data "github.com/jasondborneman/solar3/Solar3DataStorage"
	sd "github.com/jasondborneman/solar3/SolarData"
	se "github.com/jasondborneman/solar3/SolarEdge"
	"os"
	"strconv"
	"time"
)

func GetData() sd.SolarData {
	var retVal sd.SolarData

	now := time.Now()
	var siteID int
	siteID, _ = strconv.Atoi(os.Getenv("SOLAREDGE_SITEID"))
	latitude := os.Getenv("SOLAR3_LATITUDE")
	longitude := os.Getenv("SOLAR3_LONGITUDE")
	var site *se.Site
	var powerNow *se.Power
	var sunMoon *ipgl.SunMoonInfo
	var openWeather *ow.OpenWeather
	site, _ = se.GetSolarSiteInfo(siteID)
	powerNow, _ = se.GetPowerData(siteID)
	sunMoon, _ = ipgl.GetSunPosition(latitude, longitude)
	openWeather, _ = ow.GetWeather(latitude, longitude)
	retVal.Site = site
	retVal.DateTime = now
	for i := range powerNow.PowerDetails.Meters {
		if powerNow.PowerDetails.Meters[i].Type == "Production" {
			valueIndex := len(powerNow.PowerDetails.Meters[i].Values) - 1
			retVal.PowerGen = powerNow.PowerDetails.Meters[i].Values[valueIndex].Value
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

func Run() {
	var data Solar3Data
	data = GetData()
	xVals, yVals, errSave := s3data.SaveToFirestore(data)
	saved := true
	if errSave != nil {
		saved = false
	}
	graphBytes := g.CreateGraph(xVals, yVals)
	errSaveJpeg := g.SaveGraph(graphBytes, "chart")
	pngSaved := true
	if errSaveJpeg != nil {
		pngSaved = false
	}
	message := `
DateTime: %s
Last Reported Power: %.2f
Cloud Cover: %d
Sun Azimuth: %.2f
Sun Altitude: %.2f`
	message = fmt.Sprintf(
		message,
		fmt.Sprintf("%02d-%02d-%d %02d:%02d",
			data.DateTime.Month(),
			data.DateTime.Day(),
			data.DateTime.Year(),
			data.DateTime.Hour(),
			data.DateTime.Minute()),
		data.PowerGen,
		data.CloudCover,
		data.SunAzimuth,
		data.SunAltitude)
	tweetErr := tw.TweetWithMedia(message, graphBytes)
	tweeted := true
	if tweetErr != nil {
		tweeted = false
	}
	fmt.Println("----------------------------")
	fmt.Printf("Saved To Firestore?: %t\n", saved)
	fmt.Printf("Saved Graph?:        %t\n", pngSaved)
	fmt.Printf("Tweeted?:            %t\n", tweeted)
	fmt.Println(message)
}
