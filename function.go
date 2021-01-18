package solar3

import (
	"net/http"
	"os"

	s3 "github.com/jasondborneman/solar3/Solar3App"
)

func Solar3(w http.ResponseWriter, r *http.Request) {
	doTweet := os.Getenv("DO_TWEET") == "true"
	doSaveGraph := os.Getenv("DO_SAVEGRAPH") == "true"
	s3.Run(doTweet, doSaveGraph, false)
}
