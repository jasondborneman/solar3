package solar3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	s3 "github.com/jasondborneman/solar3/Solar3App"
)

type StupidAuth struct {
	stupidAuth string
}

func Solar3(w http.ResponseWriter, r *http.Request) {
	var sa StupidAuth
	err := json.NewDecoder(r.Body).Decode(&sa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stupidAuthLocal := os.Getenv("STUPID_AUTH")
	if sa.stupidAuth == stupidAuthLocal {
		doTweet := os.Getenv("DO_TWEET") == "true"
		doSaveGraph := os.Getenv("DO_SAVEGRAPH") == "true"
		s3.Run(doTweet, doSaveGraph, false)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Unauthorized access attempt. [%s | %s]\n", sa.stupidAuth, stupidAuthLocal)
		returnMessage := fmt.Sprintf("Unauthorized. [%s]", bodyString)
		fmt.Fprint(w, returnMessage)
	}
}
