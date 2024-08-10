package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Setups []Setup
}

type Setup struct {
	Name   string
	Groups []Group
}

type Group struct {
	Name     string
	Settings []Setting
}

type Setting struct {
	EC Encoder
	PB PushButton
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

func parseConfig() (cfg Config, err error) {
	err = json.NewDecoder(os.Stdin).Decode(&cfg)
	return
}
