package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

// MIDISender handles sending sysex messages via MIDI
type MIDISender struct {
	out drivers.Out
}

// NewMIDISender creates a new MIDI sender for the specified output port
func NewMIDISender(portName string) (*MIDISender, error) {
	// Get available output ports
	outs, err := drivers.Outs()
	if err != nil {
		return nil, fmt.Errorf("failed to get MIDI outputs: %w", err)
	}

	// Find the specified port
	var targetOut drivers.Out
	for _, out := range outs {
		if out.String() == portName {
			targetOut = out
			break
		}
	}

	if targetOut == nil {
		// List available ports for debugging
		fmt.Println("Available MIDI output ports:")
		for i, out := range outs {
			fmt.Printf("  %d: %s\n", i, out.String())
		}
		return nil, fmt.Errorf("MIDI output port '%s' not found", portName)
	}

	// Open the output port
	if err := targetOut.Open(); err != nil {
		return nil, fmt.Errorf("failed to open MIDI output port: %w", err)
	}

	return &MIDISender{out: targetOut}, nil
}

// SendSysex sends a sysex message to the connected device
func (ms *MIDISender) SendSysex(sysexData []byte) error {
	// Create sysex message
	msg := midi.Message(sysexData)

	// Send the message
	if err := ms.out.Send(msg); err != nil {
		return fmt.Errorf("failed to send sysex message: %w", err)
	}

	// Give some time for the message to be processed
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Close closes the MIDI output port
func (ms *MIDISender) Close() error {
	if ms.out != nil {
		return ms.out.Close()
	}
	return nil
}

// SendSysexToDevice sends a sysex message to the Faderfox EC4 device
func SendSysexToDevice(sysexData []byte) error {
	portName := "Faderfox EC4"

	sender, err := NewMIDISender(portName)
	if err != nil {
		return fmt.Errorf("failed to create MIDI sender: %w", err)
	}
	defer sender.Close()

	fmt.Printf("Sending sysex message to '%s'...\n", portName)
	fmt.Printf("Message size: %d bytes\n", len(sysexData))

	// Send the sysex message
	if err := sender.SendSysex(sysexData); err != nil {
		return fmt.Errorf("failed to send sysex: %w", err)
	}

	return nil
}
