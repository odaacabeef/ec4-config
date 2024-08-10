package main

import (
	"fmt"
	"log"
)

func main() {

	cfg, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	// TODO: construct sysex message and send to EC4 ¯\_(ツ)_/¯
}
