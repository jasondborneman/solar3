package Twitter

import (
	"fmt"
	"os"

	"github.com/dghubble/oauth1"
	"github.com/jasondborneman/go-twitter/twitter"
)

func TweetWithMedia(message string, media []byte) error {
	config := oauth1.NewConfig(os.Getenv("TWITTER_CONSUMERKEY"), os.Getenv("TWITTER_CONSUMERSECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESSTOKEN"), os.Getenv("TWITTER_ACCESSSECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	tweetParams := &twitter.StatusUpdateParams{}

	graphRes, _, err2 := client.Media.Upload(media, "image/png")
	if err2 != nil {
		fmt.Println(err2)
		return err2
	}
	tweetParams.MediaIds = []int64{}

	if graphRes.MediaID > 0 {
		tweetParams.MediaIds = append(tweetParams.MediaIds, graphRes.MediaID)
	}
	_, _, err3 := client.Statuses.Update(message, tweetParams)
	if err3 != nil {
		fmt.Println(err3)
		return err3
	}
	return nil
}
