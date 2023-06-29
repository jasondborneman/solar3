package Mastodon

import (
	"context"
	"log"
	"os"

	"github.com/jasondborneman/go-mastodon"
)

func TootWithMedia(message string, media [][]byte) error {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       "https://botsin.space",
		ClientID:     os.Getenv("MASTODON_CLIENTID"),
		ClientSecret: os.Getenv("MASTODON_CLIENTSECRET"),
		AccessToken:  os.Getenv("MASTODON_TOKEN"),
	})

	var mediaIDs []mastodon.ID
	for _, mediaBytes := range media {
		uploadRes, err := client.UploadMediaFromBytes(context.Background(), mediaBytes)
		if err != nil {
			log.Fatalf("MastoUploadMediaError: %v", err)
			return err
		}
		mediaIDs = append(mediaIDs, uploadRes.ID)
	}
	theToot := mastodon.Toot{
		Status:   message,
		MediaIDs: mediaIDs,
	}
	_, err := client.PostStatus(context.Background(), &theToot)
	if err != nil {
		log.Fatalf("MastoTootError: %v", err)
		return err
	}
	return nil
}
