package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("API KEY: %s\n", os.Getenv("RIOT_API_KEY"))
}
