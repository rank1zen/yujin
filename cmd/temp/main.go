package main

import (
	"fmt"

	"github.com/rank1zen/yujin/cmd/inter"
	"github.com/rank1zen/yujin/internal/postgresql"
)

func main() {
	q, err := inter.NewPostgreSQL()

	if err != nil {
		panic(err)
	}

	da := postgresql.NewSummonerDA(q)
	
	suu, err := da.FindRecent("zxcx")
	fmt.Printf("suu: %v\n", suu)
}
