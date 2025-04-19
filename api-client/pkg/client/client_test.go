package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	c, err := New("http://localhost:1781")
	require.NoError(t, err)
	ctx := context.Background()

	t.Logf("request id: %s", c.GenerateID())
	t.Run("Auth", func(t *testing.T) {
		login := fmt.Sprintf("testuser%d", time.Now().Unix())
		password := "password"
		_, err := c.AuthRegister(ctx, login, password)
		require.NoError(t, err)
		token, err := c.AuthLogin(ctx, login, password)
		require.NoError(t, err)
		c.SetToken(token)
		_, err = c.AuthWhoami(ctx)
		require.NoError(t, err)
	})
}
