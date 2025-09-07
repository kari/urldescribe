package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kari/urldescribe"
)

func main() {
	ctx := context.Background()
	if len(os.Args) > 1 {
		res, _ := urldescribe.DescribeURL(ctx, os.Args[1])
		fmt.Println(res)
	} else {
		fmt.Println("Usage: title <url>")
	}
}
