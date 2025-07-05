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

	fmt.Println("Configuration loaded successfully")

	// Generate sysex message
	generator := NewSysexGenerator(cfg)
	sysex, err := generator.GenerateSysex()
	if err != nil {
		log.Fatal("Failed to generate sysex:", err)
	}

	// Send sysex via MIDI
	err = SendSysexToDevice(sysex)
	if err != nil {
		log.Fatal("Failed to send sysex via MIDI:", err)
	}

	fmt.Println("Sysex message sent successfully!")
}
