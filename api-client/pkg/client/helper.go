package client

import (
	"api-client/internal/apiclient"
	"fmt"
)

func parseApiError(e apiclient.Error) error {
	return fmt.Errorf("unexpected error: %s, %s", e.Code, e.Message) //TODO
}
