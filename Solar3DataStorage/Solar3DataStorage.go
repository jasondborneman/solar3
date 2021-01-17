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

func SaveToFirestore(data sd.SolarData) ([]float64, []float64, error) {
	ctx := context.Background()
	dataClient := createClient(ctx)
	loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
	_, _, fbErr := dataClient.Collection("solar3").Add(ctx, map[string]interface{}{
		"dateNano":     time.Now().In(loc).UnixNano(),
		"SiteId":       data.Site.Details.ID,
		"DateTime":     data.DateTime,
		"PowerGen":     data.PowerGen,
		"CloudCover":   data.CloudCover,
		"Temp":         data.Temp,
		"Pressure":     data.Pressure,
		"Humidity":     data.Humidity,
		"SunAzimuth":   data.SunAzimuth,
		"SunAltitude":  data.SunAltitude,
		"SnowPastHour": data.SnowOneHr,
		"RainPastHour": data.RainOneHr,
		"WeatherID":    data.WeatherID,
	})
	if fbErr != nil {
		fbErrMessage := fmt.Sprintf("Failed adding solar data to database: %v", fbErr)
		fmt.Println(fbErrMessage)
		return nil, nil, fbErr
	}

	xVals := []float64{}
	yVals := []float64{}
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
			yVals = append(yVals, float64(doc.Data()["PowerGen"].(float64)))
		}
	}

	dataClient.Close()
	return xVals, yVals, nil
}
