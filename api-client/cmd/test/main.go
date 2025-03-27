package main

import (
	"api-client/pkg/client"
	"context"
	"fmt"
)

func main() {
	c := client.Client{}
	res, err := c.GetMemeByID(context.Background(), "0195d17ca1f47a01a4cc837f853da8f3a")
	fmt.Println(res)
	fmt.Println(err)
}
