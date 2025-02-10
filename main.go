package main

import (
	"fmt"
	"gator/internal/config"
)

func main() {
	cfg, err := config.Read()

	if err != nil {
		panic(err)
	}
	cfg.SetUser("piotr")
	cfg, err = config.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg.DbURL)

}
