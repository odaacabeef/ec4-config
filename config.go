package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Setups map[string]struct {
		Name   string
		Groups map[string]string
	}
	Groups map[string]struct {
		Name    string
		Dynamic []struct {
			Setting  string
			Behavior string
			From     int
			To       int
			Step     int
			EC       Encoder
			PB       PushButton
		}
		Static map[string]struct {
			EC Encoder
			PB PushButton
		}
	}
}

func parseConfig() (cfg Config, err error) {
	yamlFile, err := os.ReadFile("ec4.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	return
}
