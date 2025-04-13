package client

import (
	"api-client/internal/apiclient"
	"api-client/pkg/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
)

type Client struct {
	api *apiclient.ClientWithResponses
}

func New(url string) (*Client, error) {
	c, err := apiclient.NewClientWithResponses(url)
	if err != nil {
		return nil, fmt.Errorf("can't create client with responses: %w", err)
	}
	return &Client{api: c}, nil

}

// Meme

func (c *Client) PostMeme(ctx context.Context, meme models.Meme) (models.MemeID, error) {
	req := apiclient.MemeCreate{BoardId: string(meme.BoardID)}
	if meme.Descriptions != nil {
		desc := map[string]any{}
		for k, v := range meme.Descriptions {
			desc[k] = v
		}
		req.Description = &desc
	}
	if meme.Filename != "" {
		req.Filename = &meme.Filename
	}

	resp, err := c.api.PostMemeWithResponse(ctx, req)
	if err != nil {
		return models.MemeID(""), fmt.Errorf("can't post meme: %w", err)
	}

	if resp.JSON201 != nil {
		return models.MemeID(resp.JSON201.Id), nil
	}
	if resp.JSON400 != nil {
		return models.MemeID(""), parseApiError(*resp.JSON400)
	}
	return models.MemeID(""), fmt.Errorf("unexpected response")
}

func (c *Client) GetMemeByID(ctx context.Context, id models.MemeID) (models.Meme, error) {

	resp, err := c.api.GetMemeByIDWithResponse(ctx, apiclient.MemeId(id))
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't get meme: %w", err)
	}
	if resp.JSON200 != nil {
		return convertToModel(*resp.JSON200), nil
	}
	if resp.JSON404 != nil {
		return models.Meme{}, parseApiError(*resp.JSON404)
	}
	return models.Meme{}, fmt.Errorf("unexpected response")

}

func (c *Client) PutMeme(ctx context.Context, meme models.Meme) (models.Meme, error) {
	req := apiclient.MemeUpdate{}
	if meme.BoardID != "" {
		req.BoardId = (*string)(&meme.BoardID)
	}
	if meme.Filename != "" {
		req.Filename = &meme.Filename
	}
	if meme.Descriptions != nil {
		dsc := convertMapToAny(meme.Descriptions)
		req.Description = &dsc
	}

	resp, err := c.api.UpdateMemeByIDWithResponse(ctx, apiclient.MemeId(meme.ID), req)
	if err != nil {
		return models.Meme{}, fmt.Errorf("can't update meme: %w", err)
	}
	if resp.JSON200 != nil {
		return convertToModel(*resp.JSON200), nil
	}
	if resp.JSON404 != nil {
		return models.Meme{}, parseApiError(*resp.JSON404)
	}
	if resp.JSON400 != nil {
		return models.Meme{}, parseApiError(*resp.JSON400)
	}
	return models.Meme{}, fmt.Errorf("unexpected response")
}

func (c *Client) DeleteMeme(ctx context.Context, id models.MemeID) error {
	resp, err := c.api.DeleteMemeByIDWithResponse(ctx, apiclient.MemeId(id))
	if err != nil {
		return fmt.Errorf("can't delete meme: %w", err)
	}
	if resp.StatusCode() == 204 {
		return nil
	}
	if resp.JSON404 != nil {
		return parseApiError(*resp.JSON404)
	}

	return fmt.Errorf("unexpected response: %s", "")
}

// Media

func (c *Client) PutMedia(ctx context.Context, media models.Media, filename string) error {
	body, contentType, err := createMultipart("media", filename, media.Body)
	if err != nil {
		return fmt.Errorf("can't create multipart: %w", err)
	}

	resp, err := c.api.PutMediaByIDWithBodyWithResponse(ctx, apiclient.MediaId(media.ID), contentType, body)
	if err != nil {
		return fmt.Errorf("can't put media: %w", err)
	}
	if resp.StatusCode() == 200 {
		return nil
	}
	if resp.JSON400 != nil {
		return parseApiError(*resp.JSON400)
	}
	return fmt.Errorf("unexpected response: %s", resp.Body)
}

func (c *Client) GetMedia(ctx context.Context, id models.MediaID) (models.Media, error) {
	resp, err := c.api.GetMediaByIDWithResponse(ctx, apiclient.MediaId(id))
	if err != nil {
		return models.Media{}, fmt.Errorf("can't put media: %w", err)
	}
	if resp.StatusCode() == 200 {
		return models.Media{ID: id, Body: resp.Body}, nil
	}
	if resp.JSON404 != nil {
		return models.Media{}, parseApiError(*resp.JSON404)
	}
	return models.Media{}, fmt.Errorf("unexpected response")
}

// Search
func (c *Client) SearchMemeByBoardID(ctx context.Context, board_id models.MediaID, desc map[string]string) ([]models.Meme, error) {
	panic("uni")
}

func convertToModel(meme apiclient.Meme) models.Meme {
	dsc := map[string]string{}
	if meme.Description != nil {
		for k, v := range *meme.Description {
			dsc[k] = v.(string)
		}
	}
	m := models.Meme{
		ID:           models.MemeID(meme.Id),
		BoardID:      models.BoardID(meme.BoardId),
		Descriptions: dsc,
		Filename:     unptr(meme.Filename),
		CreatedAt:    meme.CreatedAt,
		UpdatedAt:    meme.UpdatedAt,
	}
	return m
}

func unptr[T any](v *T) (res T) {
	if v != nil {
		res = *v
	}
	return res
}

func convertMapToAny(o map[string]string) map[string]any {
	dsc := map[string]any{}
	for k, v := range o {
		dsc[k] = v
	}
	return dsc
}

func createMultipart(fieldName, filename string, data []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Создаем часть формы для файла
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, "", err
	}

	// Записываем данные в часть формы
	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	// Закрываем writer - это важно для корректного завершения multipart-сообщения
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}
