package main

import (
	"api-client/pkg/client"
	"context"
	"fmt"
)

func main() {
	c, err := client.New("http://localhost:1781")
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	c.SetToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDE5NjQzNjVmYWFhNzI4YzhlY2I2MTJhYzcwOTAzNDciLCJleHAiOjE3NDUxNTQ1ODIsIm5iZiI6MTc0NTA2ODE4MiwiaWF0IjoxNzQ1MDY4MTgyfQ.sVr0OkkNhLDYaGmDlYeVKbq26HcMV1C69IQ523QsMb0")
	fmt.Println("id", c.GenerateID())
	fmt.Println(c.About(ctx))
	fmt.Println(c.AuthLogin(ctx, "spr", "superstrongpassword"))
	fmt.Println(c.AuthWhoami(ctx))
	fmt.Println(c.GetMediaByID(ctx, "0196435be7e9701db5fcbabaa67dabdc"))

	// meme, err = c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a")
	// fmt.Println("meme", meme)
	// fmt.Println("err", err)

	// meme, err = c.GetMemeByID(ctx, "01962b847b4e7e41a2d78ef05bf56f7a2")
	// fmt.Println("meme", meme)
	// fmt.Println("err", err)
}
