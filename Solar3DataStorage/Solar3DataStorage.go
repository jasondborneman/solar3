package Solar3DataStorage

import (
	"context"
	"fmt"
	"os"
	"time"

	sd "github.com/jasondborneman/solar3/SolarData"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func createClient(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("GCP_PROJECT")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fsErrMessage := fmt.Sprintf("Failed to create client: %v", err)
		fmt.Println(fsErrMessage)
	}
	return client
}

func SaveToFirestore(data sd.SolarData) ([]float64, []float64, []int64, float64, error) {
	ctx := context.Background()
	dataClient := createClient(ctx)
	_, _, fbErr := dataClient.Collection("solar3").Add(ctx, map[string]interface{}{
		"dateNano":       data.PowerGen.Date.UnixNano(),
		"SiteId":         data.Site.Details.ID,
		"DateTime":       data.PowerGen.Date,
		"DateTimeStored": time.Now(),
		"PowerGen":       data.PowerGen.Value,
		"CloudCover":     data.CloudCover,
		"Temp":           data.Temp,
		"Pressure":       data.Pressure,
		"Humidity":       data.Humidity,
		"SunAzimuth":     data.SunAzimuth,
		"SunAltitude":    data.SunAltitude,
		"SnowPastHour":   data.SnowOneHr,
		"RainPastHour":   data.RainOneHr,
		"WeatherID":      data.WeatherID,
	})
	if fbErr != nil {
		fbErrMessage := fmt.Sprintf("Failed adding solar data to database: %v", fbErr)
		fmt.Println(fbErrMessage)
		return nil, nil, nil, 0, fbErr
	}

	iterMax := dataClient.Collection("solar3").OrderBy("PowerGen", firestore.Desc).StartAfter(time.Now().AddDate(0, -1, 0).UnixNano()).Limit(1).Documents(ctx)
	doc, maxPowerErr := iterMax.Next()
	if maxPowerErr != nil {
		mpErrMessage := fmt.Sprintf("Failed to get max power production: %s", maxPowerErr)
		fmt.Println(mpErrMessage)
		return nil, nil, nil, 0, maxPowerErr
	}
	maxPower := doc.Data()["PowerGen"].(float64)

	xVals := []float64{}
	powerYVals := []float64{}
	cloudYVals := []int64{}
	iter := dataClient.Collection("solar3").OrderBy("dateNano", firestore.Asc).StartAfter(time.Now().AddDate(0, -1, 0).UnixNano()).Documents(ctx)
	for {
		doc, fbReadErr := iter.Next()
		if fbReadErr == iterator.Done {
			break
		}
		if fbReadErr != nil {
			fbReadErrMessage := fmt.Sprintf("Failed retrieving solar data: %v", fbReadErr)
			fmt.Println(fbReadErrMessage)
			break
		} else {
			xVals = append(xVals, float64(doc.Data()["dateNano"].(int64)))
			powerYVals = append(powerYVals, float64(doc.Data()["PowerGen"].(float64)))
			cloudYVals = append(cloudYVals, int64(doc.Data()["CloudCover"].(int64)))
		}
	}

	dataClient.Close()
	return xVals, powerYVals, cloudYVals, maxPower, nil
}
