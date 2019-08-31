package main

import (
	"fmt"
	"os"

	"github.com/kari/urldescribe"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println(urldescribe.DescribeURL(os.Args[1]))
	} else {
		fmt.Println("Usage: title <url>")
	}
}
