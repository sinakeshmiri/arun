package main

import (
	"fmt"

	"github.com/sinakeshmiri/arun/packages/wrapper"
)

func main() {
	err := wrapper.Make("./t.zip")
	if err != nil {
		fmt.Println(err)
	}
}
