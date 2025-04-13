package client

import (
	"context"
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := New("http://localhost:1781")
	if err != nil {
		fmt.Println("err", err)
	}

	ctx := context.Background()

	meme, err := c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a")
	fmt.Println("meme", meme)
	fmt.Println("err", err)

	meme, err = c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a2")
	fmt.Println("meme", meme)
	fmt.Println("err", err)

}
