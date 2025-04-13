package main

import (
	"api-client/pkg/client"
	"api-client/pkg/models"
	"context"
	"fmt"
)

func main() {
	c, err := client.New("http://localhost:1781")
	if err != nil {
		fmt.Println("err", err)
	}

	ctx := context.Background()
	meme := models.Meme{
		BoardID:      "",
		Filename:     "file.png",
		Descriptions: map[string]string{"general": "a"},
	}
	id, err := c.PostMeme(ctx, meme)
	fmt.Println(id, err)
	// meme, err = c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a")
	// fmt.Println("meme", meme)
	// fmt.Println("err", err)

	// meme, err = c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a2")
	// fmt.Println("meme", meme)
	// fmt.Println("err", err)
}
