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

func CreateGraph(xVals, powerYVals, yVals []float64, maxPower float64, yName string) []byte {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: func(v interface{}) string {
				loc, _ := time.LoadLocation("America/Indiana/Indianapolis")
				typed := v.(float64)
				typedDate := time.Unix(0, int64(typed)).In(loc)
				return fmt.Sprintf("%02d:%02d", typedDate.Hour(), typedDate.Minute())
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
			Range: &chart.ContinuousRange{
				Min: 0.0,
				Max: maxPower + 1000.0,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Power Generated",
				XValues: xVals,
				YValues: powerYVals,
			},
			chart.ContinuousSeries{
				Name:    yName,
				YAxis:   chart.YAxisSecondary,
				XValues: xVals,
				YValues: yVals,
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
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
