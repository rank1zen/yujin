package main

import (
	"fmt"
	"log"

	"github.com/rank1zen/yujin/cmd/internal"
	"github.com/rank1zen/yujin/internal/rest"
)

func main() {
	r := rest.InitRouter()

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}

	db, err := internal.NewPostgreSQL()

	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("db: %v\n", db)
}

