package Mastodon

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-mastodon"
)

func TootWithMedia(message string, media []byte) error {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       "https://botsin.space",
		ClientID:     os.Getenv("MASTODON_CLIENTID"),
		ClientSecret: os.Getenv("MASTODON_CLIENTSECRET"),
	})
	err := client.Authenticate(context.Background(), os.Getenv("MASTODON_USER"), os.Getenv("MASTODON_PWD"))
	if err != nil {
		log.Fatal(fmt.Sprintf("MastoAuthError: %v", err))
		return err
	}

	uploadRes, err := client.UploadMediaFromBytes(context.Background(), media)
	if err != nil {
		log.Fatal(fmt.Sprintf("MastoUploadMediaError: %v", err))
		return err
	}
	var mediaIDs []mastodon.ID
	mediaIDs[0] = uploadRes.ID
	theToot := mastodon.Toot{
		Status:   message,
		MediaIDs: mediaIDs,
	}
	_, err = client.PostStatus(context.Background(), &theToot)
	if err != nil {
		log.Fatal(fmt.Sprintf("MastoUploadMediaError: %v", err))
		return err
	}
	return nil
}
