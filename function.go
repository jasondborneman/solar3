package solar3

import (
	"fmt"
	"net/http"
	"os"

	s3 "github.com/jasondborneman/solar3/Solar3App"
)

func Solar3(w http.ResponseWriter, r *http.Request) {
	stupidAuth := r.Header.Get("Stupid-Auth")
	stupidAuthLocal := os.Getenv("STUPID_AUTH")
	if stupidAuth == stupidAuthLocal {
		doTweet := os.Getenv("DO_TWEET") == "true"
		doSaveGraph := os.Getenv("DO_SAVEGRAPH") == "true"
		s3.Run(doTweet, doSaveGraph, false)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Unauthorized access attempt")
		fmt.Fprint(w, "Unauthorized")
	}
}
