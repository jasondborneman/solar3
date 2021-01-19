package solar3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	s3 "github.com/jasondborneman/solar3/Solar3App"
)

func Solar3(w http.ResponseWriter, r *http.Request) {
	var d struct {
		StupidAuth string `json:"stupidAuth"`
	}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stupidAuthLocal := os.Getenv("STUPID_AUTH")
	if d.StupidAuth == stupidAuthLocal {
		doTweet := os.Getenv("DO_TWEET") == "true"
		doSaveGraph := os.Getenv("DO_SAVEGRAPH") == "true"
		s3.Run(doTweet, doSaveGraph, false)
		fmt.Fprint(w, "Success")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
	}
}
