package client

import (
	"api-client/internal/requester"
	"api-client/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Url string
}

// Meme

func (c *Client) PostMeme(ctx context.Context, meme models.Meme) (models.MemeID, error) {
	req := requester.Request{
		Method: "POST",
		Url:    fmt.Sprintf("%s/memes", c.Url),
		Body:   meme,
	}
	res := struct {
		ID models.MemeID `json:"id"`
	}{}

	err := processAndParse(req, &res)
	return res.ID, err

}

func (c *Client) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {
	req := requester.Request{
		Method: "GET",
		Url:    fmt.Sprintf("%s/memes/%s", c.Url, id),
	}
	var res models.Meme
	err := processAndParse(req, &res)
	return res, err
}

func (c *Client) PutMeme(ctx context.Context, meme models.Meme) error {
	req := requester.Request{
		Method: "PUT",
		Url:    fmt.Sprintf("%s/memes/%s", c.Url, meme.ID),
		Body:   meme,
	}
	res := struct{}{}
	err := processAndParse(req, &res)
	return err
}

func (c *Client) DeleteMeme(ctx context.Context, id models.MemeID) error {
	req := requester.Request{
		Method: "DELETE",
		Url:    fmt.Sprintf("%s/memes/%s", c.Url, id),
	}
	res := struct{}{}
	err := processAndParse(req, &res)
	return err
}

// // Media

func (c *Client) PutMedia(ctx context.Context, media models.Media, filename string) error {
	req := requester.Request{
		Method: "PUT",
		Url:    fmt.Sprintf("%s/media/%s", c.Url, media.ID),
		MultipartFiles: map[string]struct {
			Data io.Reader
			Name string
		}{
			"media": {
				Data: bytes.NewBuffer(media.Body),
				Name: filename,
			},
		},
	}
	res := struct{}{}
	err := processAndParse(req, &res)
	return err

}

func (c *Client) GetMedia(ctx context.Context, id models.MediaID) (models.Media, error) {
	req := requester.Request{
		Method: "PUT",
		Url:    fmt.Sprintf("%s/media/%s", c.Url, id),
	}
	resp, err := req.Do()
	if err != nil {
		return models.Media{}, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()
	res := models.Media{ID: id}
	res.Body, err = io.ReadAll(resp.Body)
	if err != nil {
		return models.Media{}, fmt.Errorf("can't read body: %w", err)
	}
	return res, nil
}

func makeJSONReqest(method string, url string, body any) (int, io.Reader, error) {
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("can't marshal body: %w", err)
		}
	}
	bodyReader := bytes.NewBuffer(bodyBytes)
	request, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, nil, fmt.Errorf("can't do request: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("can't read body: %w", err)
	}

	return resp.StatusCode, bytes.NewBuffer(bodyBytes), nil
}
