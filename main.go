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
	fmt.Printf("Found %d setups\n", len(cfg.Setups))

	// Generate sysex message
	generator := NewSysexGenerator(cfg)
	sysex, err := generator.GenerateSysex()
	if err != nil {
		log.Fatal("Failed to generate sysex:", err)
	}

	// Print sysex message
	PrintSysexHex(sysex)

	// Send sysex via MIDI
	err = SendSysexToDevice(sysex)
	if err != nil {
		log.Fatal("Failed to send sysex via MIDI:", err)
	}

	fmt.Println("Configuration sent to Faderfox EC4 successfully!")
}
