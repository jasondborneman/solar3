package solar3

import (
	"net/http"

	s3 "github.com/jasondborneman/solar3/Solar3App"
)

func Solar3(w http.ResponseWriter, r *http.Request) {
	s3.Run()
}
