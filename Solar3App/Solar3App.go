package Solar3App

import (
	ipgl "jborneman/solar3/IPGeoLocation"
	ow "jborneman/solar3/OpenWeather"
	se "jborneman/solar3/SolarEdge"
	"os"
	"strconv"
	"time"
)

type Solar3Data struct {
	Site        *se.Site
	DateTime    time.Time
	PowerGen    float64
	CloudCover  int
	Temp        float64
	Pressure    int
	Humidity    int
	SunAzimuth  float64
	SunAltitude float64
	SnowOneHr   float64
	RainOneHr   float64
	WeatherID   int
}

func GetData() Solar3Data {
	var retVal Solar3Data

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
