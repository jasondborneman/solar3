package main

import (
	"fmt"
	g "github.com/jasondborneman/solar3/Graphing"
	s3 "github.com/jasondborneman/solar3/Solar3App"
	s3data "github.com/jasondborneman/solar3/Solar3DataStorage"
	tw "github.com/jasondborneman/solar3/Twitter"
)

func main() {
	var data s3.Solar3Data
	data = s3.GetData()
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
