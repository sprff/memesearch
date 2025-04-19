package client

import (
	"api-client/internal/apiclient"
	"api-client/internal/utils"
	"api-client/pkg/models"
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func parseApiError(e apiclient.Error) error {
	if e.Id == "INVALID_REQUEST" && e.Body != nil {
		for k, v := range *e.Body {
			return models.ErrInvalidInput{Param: k, Reason: v.(string)}
		}
	}
	return fmt.Errorf("unexpected error: %s, %s (%v)", e.Id, e.Message, e.Body) //TODO
}

func convertMemeToModel(m apiclient.Meme) models.Meme {
	dsc := map[string]string{}
	for k, v := range m.Description {
		dsc[k] = v.(string)
	}
	return models.Meme{
		ID:           models.MemeID(m.Id),
		BoardID:      models.BoardID(m.BoardId),
		Descriptions: dsc,
		Filename:     m.Filename,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func convertUserToModel(u apiclient.User) models.User {
	return models.User{ID: models.UserID(u.Id), Login: u.Login}
}
func convertBoardToModel(b apiclient.Board) models.Board {
	return models.Board{
		ID:    models.BoardID(b.Id),
		Owner: models.UserID(b.Owner),
		Name:  b.Name,
	}
}

func convertScoredToModel(m apiclient.ScoredMeme) models.ScoredMeme {
	return models.ScoredMeme{
		Score: m.Score,
		Meme:  convertMemeToModel(m.Meme),
	}
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

	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(part, bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func (c Client) editorAuth() apiclient.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
		return nil

	}
}

func (c Client) editorRequestID() apiclient.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Add("X-Request-ID", c.requestId)
		return nil
	}
}

func (c Client) middlewares() []apiclient.RequestEditorFn {
	return []apiclient.RequestEditorFn{
		c.editorAuth(),
		c.editorRequestID(),
	}
}

func (c *Client) GenerateID() string {
	c.requestId = utils.GenereateUUIDv7()
	return c.requestId
}

func (c *Client) GetID() string {
	c.requestId = utils.GenereateUUIDv7()
	return c.requestId
}

func (c *Client) SetToken(token string) {
	c.token = token
}

func (c *Client) GetToken() string {
	return c.token
}

func ptr[T any](x T) *T {
	return &x
}
