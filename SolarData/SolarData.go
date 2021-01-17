package SolarData

import (
	"time"

	se "github.com/jasondborneman/solar3/SolarEdge"
)

type SolarData struct {
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
