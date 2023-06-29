package Solar3App

import (
	"fmt"
	"os"
	"strconv"
	"time"

	aq "github.com/jasondborneman/solar3/AirNow"
	g "github.com/jasondborneman/solar3/Graphing"
	ipgl "github.com/jasondborneman/solar3/IPGeoLocation"
	ma "github.com/jasondborneman/solar3/Mastodon"
	ow "github.com/jasondborneman/solar3/OpenWeather"
	s3data "github.com/jasondborneman/solar3/Solar3DataStorage"
	sd "github.com/jasondborneman/solar3/SolarData"
	se "github.com/jasondborneman/solar3/SolarEdge"
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
	var aqiData *aq.AQIData
	fmt.Println("GetSolarSiteInfo")
	site, _ = se.GetSolarSiteInfo(siteID)
	fmt.Println("GetLatestPowerData")
	powerNow, _ = se.GetLatestPowerData(siteID)
	fmt.Println("GetSunPosition")
	sunMoon, _ = ipgl.GetSunPosition(latitude, longitude)
	fmt.Println("GetWeather")
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
	fmt.Println("GetAQI")
	aqiData, _ = aq.GetAQI(latitude, longitude)
	retVal.CloudCover = openWeather.Clouds.All
	retVal.Temp = openWeather.Main.Temp
	retVal.Pressure = openWeather.Main.Pressure
	retVal.Humidity = openWeather.Main.Humidity
	retVal.SunAzimuth = sunMoon.SunAzimuth
	retVal.SunAltitude = sunMoon.SunAltitude
	retVal.SnowOneHr = openWeather.Snow.OneH
	retVal.RainOneHr = openWeather.Rain.OneH
	retVal.WeatherID = openWeather.Weather[0].ID
	retVal.PM2_5AQI = aqiData.PM2_5AQI
	retVal.PM10AQI = aqiData.PM10AQI
	return retVal
}

func Run(doToot bool, doSaveGraph bool, fixDodgyDataOnly bool) {
	loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
	saved := false
	pngSaved := false
	tooted := false
	message := ""
	if !fixDodgyDataOnly {
		var data sd.SolarData
		data = GetData()
		fmt.Println("SaveToFirestore")
		xVals, powerYVals, sunAltVals, pm2_5Vals, pm10Vals, maxPower, errSave := s3data.SaveToFirestore(data)
		saved = true
		if errSave != nil {
			saved = false
		}
		fmt.Println("CreateGraph")
		graphBytesSunAlt := g.CreateGraph(xVals, powerYVals, sunAltVals, maxPower, "Sun Altitude")
		if doSaveGraph {
			fmt.Println("SaveGraph (Sun Altitude)")
			errSaveJpeg := g.SaveGraph(graphBytesSunAlt, "chartSunAlt")
			pngSaved = true
			if errSaveJpeg != nil {
				pngSaved = false
			}
		}
		graphBytesPM2_5 := g.CreateGraph(xVals, powerYVals, pm2_5Vals, maxPower, "AQI (PM2.5)")
		if doSaveGraph {
			fmt.Println("SaveGraph (PM2.5)")
			errSaveJpeg := g.SaveGraph(graphBytesSunAlt, "chartSunPM2_5")
			pngSaved = true
			if errSaveJpeg != nil {
				pngSaved = false
			}
		}
		graphBytesPM10 := g.CreateGraph(xVals, powerYVals, pm10Vals, maxPower, "AQI (PM10)")
		if doSaveGraph {
			fmt.Println("SaveGraph (PM2.5)")
			errSaveJpeg := g.SaveGraph(graphBytesSunAlt, "chartSunPM210")
			pngSaved = true
			if errSaveJpeg != nil {
				pngSaved = false
			}
		}
		images := [][]byte{graphBytesSunAlt, graphBytesPM2_5, graphBytesPM10}
		message = `
DateTime: %s
Last Reported Power: %.2f
Cloud Cover: %d
Sun Altitude: %.2f
AQI (PM2.5): %d
AQI (PM10): %d`
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
			data.SunAltitude,
			data.PM2_5AQI,
			data.PM10AQI)
		tooted = false
		if doToot {
			if data.SunAltitude > 0 {
				fmt.Println("TootWithMedia")
				tootErr := ma.TootWithMedia(message, images)
				tooted = true
				if tootErr != nil {
					tooted = false
				}
			} else {
				fmt.Printf("It's night, no point Tooting!")
			}
		}
	}
	fmt.Println("FixingDodgyData")
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
		if fixedPower.Value > 0 {
			updateErr := s3data.UpdatePowerDataAt(fixedPower)
			if updateErr != nil {
				fmt.Printf("Error updating dodgy power data at [%s]: %s", t.In(loc), updateErr)
			}
			fixedCount++
		}
	}
	fmt.Println("----------------------------")
	if !fixDodgyDataOnly {
		fmt.Printf("Saved To Firestore?: %t\n", saved)
		fmt.Printf("Saved Graph?:        %t\n", pngSaved)
		fmt.Printf("Tooted?:            %t\n", tooted)
		fmt.Println(message)
	}
	fmt.Printf("Dodgy Data Count:    %d\n", len(dodgyTimes))
	fmt.Printf("Fixed Data Count:    %d\n", fixedCount)
	fmt.Printf("Fix Data Err Count:  %d\n", fixErrCount)

}
