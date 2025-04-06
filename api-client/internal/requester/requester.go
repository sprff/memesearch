package requester

import (
	"api-client/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Request struct {
	Method         string
	Url            string
	Body           any
	MultipartFiles map[string]struct {
		Data io.Reader
		Name string
	}
}

func (r *Request) Do() (*http.Response, error) {
	var requestBody bytes.Buffer

	if r.Body != nil && r.MultipartFiles != nil {
		panic("either request.Body or request.MultipartFiles must be not nil")
	}

	if r.Body != nil {

		bodyBytes, err := json.Marshal(r.Body)
		if err != nil {
			return nil, fmt.Errorf("can't marshal body: %w", err)
		}
		_, err = requestBody.Write(bodyBytes)
		if err != nil {
			return nil, fmt.Errorf("can't write to buffer: %w", err)
		}
	}

	contentType := ""
	if r.MultipartFiles != nil {
		writer := multipart.NewWriter(&requestBody)
		for field, file := range r.MultipartFiles {
			part, err := writer.CreateFormFile(field, file.Name)
			if err != nil {
				return nil, fmt.Errorf("can't create file part: %w", err)
			}
			_, err = io.Copy(part, file.Data)
			if err != nil {
				return nil, fmt.Errorf("can't copy file data: %w", err)
			}
		}
		writer.Close()
		contentType = writer.FormDataContentType()
	}

	request, err := http.NewRequest(r.Method, r.Url, &requestBody)
	request.Header.Set("Content-Type", contentType)
	if err != nil {
		return nil, fmt.Errorf("can't create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	return resp, nil
}

func ReadResponse[T any](r io.Reader, input *T, apiErr *error) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("can't read: %w", err)
	}

	res := struct {
		Status  string `json:"status"`
		Data    any    `json:"data"`
		ErrData any    `json:"err_data"`
	}{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("can't unmarshal response: %w", err)
	}

	if res.Status != "OK" {
		*apiErr = ParseError(res.Status, res.ErrData)
		return nil
	}
	m, err := json.Marshal(res.Data)
	if err != nil {
		return fmt.Errorf("can't marshal res.Data: %w", err)
	}
	err = json.Unmarshal(m, &input)
	if err != nil {
		return fmt.Errorf("can't unmarshal input: %w", err)
	}
	return nil
}

func ParseError(status string, data any) error {
	switch status {
	case "BOARD_NOT_FOUND":
		return models.ErrBoardNotFound
	case "MEDIA_NOT_FOUND":
		return models.ErrMediaNotFound
	case "MEDIA_IS_REQUIRED":
		return models.ErrMediaIsRequired
	case "MEME_NOT_FOUND":
		return models.ErrMemeNotFound
	case "INVALID_INPUT":
		return remarshalError[models.ErrInvalidInput](data)
	}
	panic(fmt.Sprintf("Unhandled error occured, status: %s", status))
}

func remarshalError[T error](data any) error {
	var res T
	m, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("can't marshal data: %w", err)
	}
	err = json.Unmarshal(m, &res)
	if err != nil {
		return fmt.Errorf("can't unmarshal res: %w", err)
	}
	return res
}
