package solar3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	fmt.Println(d)
	if d.StupidAuth == stupidAuthLocal {
		doTweet := os.Getenv("DO_TWEET") == "true"
		doSaveGraph := os.Getenv("DO_SAVEGRAPH") == "true"
		s3.Run(doTweet, doSaveGraph, false)
		fmt.Fprint(w, "Success")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		fmt.Printf("Unauthorized access attempt. [%s | %s]\n", d.StupidAuth, stupidAuthLocal)
		fmt.Printf("%s", bodyString)
		returnMessage := fmt.Sprintf("Unauthorized. [%s]", bodyString)
		fmt.Fprint(w, returnMessage)
	}
}
