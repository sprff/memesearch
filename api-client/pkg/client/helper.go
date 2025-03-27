package client

import (
	"api-client/internal/requester"
	"fmt"
)

func processAndParse[T any](req requester.Request, res *T) error {
	resp, err := req.Do()
	if err != nil {
		return fmt.Errorf("can't do request: ")
	}
	var apiErr error
	err = requester.ReadResponse(resp.Body, res, &apiErr)
	if err != nil {
		return fmt.Errorf("can't read response: %w", err)
	}
	return apiErr

}
