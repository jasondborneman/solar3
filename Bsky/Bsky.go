package Bsky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var bskyUri = "https://bsky.social"

type BskyAuth struct {
	AccessJwt string `json:"accessJwt"`
}

type BskyAuthPost struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type BskyImageWithAlt struct {
	Alt   string        `json:"alt"`
	Image BskyImageBlob `json:"image"`
}

type BskyMediaPost struct {
	Type      string    `json:"$type"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Embed     struct {
		Type   string             `json:"$type"`
		Images []BskyImageWithAlt `json:"images"`
	} `json:"embed"`
}

type BskyPost struct {
	Collection string        `json:"collection"`
	Repo       string        `json:"repo"`
	Record     BskyMediaPost `json:"record"`
}

type BskyImageBlob struct {
	Type string `json:"$type"`
	Ref  struct {
		Link string `json:"$link"`
	} `json:"ref"`
	MimeType string `json:"mimeType"`
	Size     int    `json:"size"`
}

type BskyImageUploadResp struct {
	Blob BskyImageBlob `json:"blob"`
}

func PostWithMedia(message string, media [][]byte) error {
	var bskyClient = &http.Client{
		Timeout: time.Second * 10,
	}
	handle := os.Getenv("BSKY_USER")
	pass := os.Getenv("BSKY_PASS")
	bskyAuthPost := &BskyAuthPost{
		Identifier: handle,
		Password:   pass,
	}
	var authBuf bytes.Buffer
	encodeErr := json.NewEncoder(&authBuf).Encode(bskyAuthPost)
	if encodeErr != nil {
		log.Fatalf("Error encoding Bsky Auth Post: %s", encodeErr)
		return encodeErr
	}
	url := fmt.Sprintf("%s/xrpc/com.atproto.server.createSession", bskyUri)
	resp, authErr := bskyClient.Post(url, "application/json", &authBuf)
	if authErr != nil {
		log.Fatalf("Error authenticating to Bsky: %s", authErr)
		return authErr
	}
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("Non-200 Status Code Returned: %d [%s]", resp.StatusCode, url)
		log.Fatal(message)
		return errors.New(message)
	}
	defer resp.Body.Close()
	bskyAuth := &BskyAuth{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&bskyAuth)
	if decodeErr != nil {
		log.Fatalf("Error decoding Bsky Auth response: %s", decodeErr)
		return decodeErr
	}

	bskyMediaPost := &BskyMediaPost{}
	bskyMediaPost.Type = "app.bsky.feed.post"
	bskyMediaPost.Text = message
	bskyMediaPost.CreatedAt = time.Now()
	bskyMediaPost.Embed.Type = "app.bsky.embed.image"
	bskyMediaPost.Embed.Images = []BskyImageWithAlt{}
	for i, mediaBytes := range media {
		if len(mediaBytes) > 1000000 {
			log.Fatalf("Image %d is too large: %d", i, len(mediaBytes))
			return errors.New("Image too large")
		}
		url = fmt.Sprintf("%s/xrpc/com.atproto.repo.uploadBlob", bskyUri)
		uploadImgReq, uploadImgErr := http.NewRequest("POST", url, bytes.NewReader(mediaBytes))
		if uploadImgErr != nil {
			log.Fatalf("Error creating Bsky Image Upload request: %s", uploadImgErr)
			return uploadImgErr
		}
		uploadImgReq.Header.Set("Content-Type", "image/png")
		uploadImgReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bskyAuth.AccessJwt))
		uploadImgResp, uploadImgErr := bskyClient.Do(uploadImgReq)
		if uploadImgErr != nil {
			log.Fatalf("Error uploading image to Bsky: %s", uploadImgErr)
			return uploadImgErr
		}
		bskyBlob := &BskyImageUploadResp{}
		decodeErr := json.NewDecoder(uploadImgResp.Body).Decode(bskyBlob)
		if decodeErr != nil {
			log.Fatalf("Error decoding Bsky Blob response: %s", decodeErr)
			return decodeErr
		}
		bskyMediaPost.Embed.Images = append(bskyMediaPost.Embed.Images, BskyImageWithAlt{
			Alt:   fmt.Sprintf("Image %d", i),
			Image: BskyImageBlob{},
		})
		log.Printf("Bsky Blob: %#v", bskyBlob)
		log.Printf("Bsky Blob Ref: %#v", bskyBlob.Blob.Ref)
		log.Printf("Bsky Blob Ref Link: %#v", bskyBlob.Blob.Ref.Link)
		log.Printf("Bsky Blob MimeType: %#v", bskyBlob.Blob.MimeType)
		log.Printf("Bsky Blob Size: %#v", bskyBlob.Blob.Size)
		bskyMediaPost.Embed.Images[i].Image.MimeType = bskyBlob.Blob.MimeType
		bskyMediaPost.Embed.Images[i].Image.Size = bskyBlob.Blob.Size
		bskyMediaPost.Embed.Images[i].Image.Ref.Link = bskyBlob.Blob.Ref.Link
		bskyMediaPost.Embed.Images[i].Image.Type = "blob"

	}

	bskyPost := &BskyPost{
		Collection: "app.bsky.feed.post",
		Repo:       handle,
		Record:     *bskyMediaPost,
	}

	var buf bytes.Buffer
	encodeErr = json.NewEncoder(&buf).Encode(bskyPost)
	postBody := buf.String()
	log.Printf("Posting to Bsky: %s", postBody)
	log.Printf("Posting to Bsky (struct): %#v", bskyPost)
	if encodeErr != nil {
		log.Fatalf("Error encoding Bsky Media Post: %s", encodeErr)
		return encodeErr
	}
	url = fmt.Sprintf("%s/xrpc/com.atproto.repo.createRecord", bskyUri)
	postReq, postErr := http.NewRequest("POST", url, &buf)
	if postErr != nil {
		log.Fatalf("Error creating Bsky Media Post request: %s", postErr)
		return postErr
	}
	postReq.Header.Set("Content-Type", "application/json")
	postReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bskyAuth.AccessJwt))
	postResp, postErr := bskyClient.Do(postReq)
	if postErr != nil {
		log.Fatalf("Error posting Bsky Media Post: %s", postErr)
		return postErr
	}
	if postResp.StatusCode != 200 {
		b, err := io.ReadAll(postResp.Body)
		if err != nil {
			log.Fatalf("Error reading Bsky Media Post response: %s", err)
			return err
		}

		message := fmt.Sprintf("Non-200 Status Code Returned making post: %d [%s] [%s]", postResp.StatusCode, url, string(b))
		log.Fatal(message)
		return errors.New(message)
	}
	return nil
}
