package main

import (
	"fmt"
	"os"

	"github.com/mr-chelyshkin/NetSurf/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
}
