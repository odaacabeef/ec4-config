package main

import (
	"fmt"
	"log"
)

type Setup struct {
	Name   string
	Groups map[string]Group
}

type Group struct {
	Name        string
	Encoders    map[string]Encoder
	PushButtons map[string]PushButton
}

type Encoder struct {
	Channel int
	Number  int
	Lower   int
	Upper   int
	Display string
	Type    string
	Mode    string
}

type PushButton struct {
	Channel int
	Number  int
	Lower   int
	Upper   int
	Display string
	Type    string
	Mode    string
}

func main() {

	cfg, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)

	// TODO: construct sysex message and send to EC4 ¯\_(ツ)_/¯
}
