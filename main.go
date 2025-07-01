package main

import (
	"fmt"
	"log"
	"os"
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

	// Write sysex to file
	err = os.WriteFile("ec4-config.syx", sysex, 0644)
	if err != nil {
		log.Fatal("Failed to write sysex file:", err)
	}

	fmt.Println("Sysex message written to ec4-config.syx")
}
