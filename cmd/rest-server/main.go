package main

import (
	"log"

	"github.com/rank1zen/yujin/internal/rest"
)

func main() {
	r := rest.InitRouter()
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}

