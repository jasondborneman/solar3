package graphing

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

func CreateGraph(xVals []float64, powerYVals []float64, cloudYVals []int64) []byte {
	cloudYValsFloat := []float64{}
	for i := range cloudYVals {
		floatVal := float64(cloudYVals[i])
		cloudYValsFloat = append(cloudYValsFloat, floatVal)
	}
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
				typed := v.(float64)
				typedDate := time.Unix(0, int64(typed)).In(loc)
				return fmt.Sprintf("%02d/%02d/%d %02d:%02d", typedDate.Month(), typedDate.Day(), typedDate.Year(), typedDate.Hour(), typedDate.Minute())
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Range: &chart.ContinuousRange{
				Min: 0.0,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xVals,
				YValues: powerYVals,
			},
			chart.ContinuousSeries{
				YAxis:   chart.YAxisSecondary,
				XValues: xVals,
				YValues: cloudYValsFloat,
			},
		},
	}

	graphBuffer := bytes.NewBuffer([]byte{})
	gErr := graph.Render(chart.PNG, graphBuffer)
	if gErr != nil {
		fmt.Println(gErr)
	}
	return graphBuffer.Bytes()
}

func SaveGraph(byteData []byte, name string) error {
	img, _, err := image.Decode(bytes.NewReader(byteData))
	if err != nil {
		log.Fatalln(err)
		return err
	}
	filename := fmt.Sprintf("./%s.png", name)
	out, _ := os.Create(filename)
	defer out.Close()

	err = png.Encode(out, img)
	if err != nil {
		return err
	}
	out.Close()
	return nil
}
