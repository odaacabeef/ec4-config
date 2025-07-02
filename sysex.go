package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// Sysex constants
const (
	SYSEX_START        = 0xF0
	SYSEX_STOP         = 0xF7
	CMD_DOWNLOAD_START = 0x41
	CMD_DOWNLOAD_TYPE  = 0x42
	CMD_APP_ID_H       = 0x43
	CMD_APP_ID_L       = 0x44
	CMD_PAGE_NUM_H     = 0x49
	CMD_PAGE_NUM_L     = 0x4A
	CMD_PAGE_CRC_H     = 0x4B
	CMD_PAGE_CRC_L     = 0x4C
	CMD_PAGE_DATA      = 0x4D
	CMD_DOWNLOAD_STOP  = 0x4F

	DEVICE_ID_EC4 = 0x0B // 11 in decimal

	// Memory addresses for EC4
	EEPROM_ADR_SETUP_KEY1  = 0x0B00
	EEPROM_ADR_SETUP_NAMES = 0x1BC0
	EEPROM_ADR_GROUP_NAMES = 0x1C00
	EEPROM_ADR_SETUP_DATA  = 0x2000
	EEPROM_ADR_SETUP_KEY2  = 0xE000
)

// SysexPage represents a page of sysex data
type SysexPage struct {
	Address uint16
	Data    []byte
}

// SysexGenerator handles the generation of sysex messages
type SysexGenerator struct {
	config Config
}

// NewSysexGenerator creates a new sysex generator
func NewSysexGenerator(cfg Config) *SysexGenerator {
	return &SysexGenerator{config: cfg}
}

// GenerateSysex creates the complete sysex message
func (sg *SysexGenerator) GenerateSysex() ([]byte, error) {
	var sysex []byte

	// Start sysex
	sysex = append(sysex, SYSEX_START, 0x00, 0x00, 0x00)

	// Download start
	sysex = append(sysex, encodeCommand(CMD_DOWNLOAD_START, DEVICE_ID_EC4)...)

	// Download type (03 = all setups)
	sysex = append(sysex, encodeCommand(CMD_DOWNLOAD_TYPE, 0x03)...)

	// App ID (firmware version)
	sysex = append(sysex, encodeCommand(CMD_APP_ID_H, 0x02)...) // Version 2.00
	sysex = append(sysex, encodeCommand(CMD_APP_ID_L, 0x00)...)

	// Generate all data areas
	pages := sg.generateAllPages()

	// Add all pages
	for _, page := range pages {
		sysex = append(sysex, sg.encodePage(page)...)
	}

	// Download stop
	sysex = append(sysex, encodeCommand(CMD_DOWNLOAD_STOP, DEVICE_ID_EC4)...)

	// End sysex
	sysex = append(sysex, SYSEX_STOP)

	return sysex, nil
}

// generateAllPages generates all the required pages for the EC4
func (sg *SysexGenerator) generateAllPages() []SysexPage {
	var pages []SysexPage

	// Generate SETUP_KEY1 pages (64 pages)
	setupKey1Pages := sg.generateSetupKey1Pages()
	pages = append(pages, setupKey1Pages...)

	// Generate SETUP_NAMES pages (1 page)
	setupNamesPages := sg.generateSetupNamesPages()
	pages = append(pages, setupNamesPages...)

	// Generate GROUP_NAMES pages (16 pages)
	groupNamesPages := sg.generateGroupNamesPages()
	pages = append(pages, groupNamesPages...)

	// Generate SETUP_DATA pages (768 pages)
	setupDataPages := sg.generateSetupDataPages()
	pages = append(pages, setupDataPages...)

	// Generate SETUP_KEY2 pages (128 pages)
	setupKey2Pages := sg.generateSetupKey2Pages()
	pages = append(pages, setupKey2Pages...)

	return pages
}

// generateSetupKey1Pages generates the SETUP_KEY1 pages
func (sg *SysexGenerator) generateSetupKey1Pages() []SysexPage {
	var pages []SysexPage
	address := uint16(EEPROM_ADR_SETUP_KEY1)

	// Total of 64 pages for all setups
	// Each setup has 4 pages (256 bytes), each group uses 0.25 pages (16 bytes)
	for setupIdx := 0; setupIdx < 16; setupIdx++ {
		for pageInSetup := 0; pageInSetup < 4; pageInSetup++ {
			pageData := make([]byte, 64)

			// Each page contains data for 4 groups
			for groupInPage := 0; groupInPage < 4; groupInPage++ {
				groupIdx := pageInSetup*4 + groupInPage
				if groupIdx >= 16 {
					break
				}

				// Get the actual group data if available
				var group *Group
				if setupIdx < len(sg.config.Setups) && groupIdx < len(sg.config.Setups[setupIdx].Groups) {
					group = &sg.config.Setups[setupIdx].Groups[groupIdx]
				}

				// Generate push button mode + command data for this group
				groupData := sg.generatePushButtonKey1Data(group)

				// Copy group data to page (16 bytes per group)
				copy(pageData[groupInPage*16:], groupData)
			}

			pages = append(pages, SysexPage{
				Address: address,
				Data:    pageData,
			})
			address += 64
		}
	}

	return pages
}

// generatePushButtonKey1Data generates push button mode + command data for a group
func (sg *SysexGenerator) generatePushButtonKey1Data(group *Group) []byte {
	data := make([]byte, 16)

	if group == nil {
		return data
	}

	// Each byte: bit 7 = mode, bits 0-6 = command number
	for i := 0; i < 16; i++ {
		if i < len(group.Settings) {
			setting := group.Settings[i]

			// Convert push button mode
			mode := sg.convertPushButtonTypeToMode(setting.PB.Mode)

			// Command number (0-127)
			command := setting.PB.Number
			if command > 127 {
				command = 127
			}

			// Combine mode and command
			data[i] = byte((mode << 7) | command)
		}
	}

	return data
}

// generateSetupNamesPages generates the SETUP_NAMES pages
func (sg *SysexGenerator) generateSetupNamesPages() []SysexPage {
	pageData := make([]byte, 64)

	// Each setup name is 4 characters
	for setupIdx := 0; setupIdx < 16; setupIdx++ {
		var name string
		if setupIdx < len(sg.config.Setups) {
			name = sg.config.Setups[setupIdx].Name
		}

		// Pad or truncate to 4 characters
		if len(name) > 4 {
			name = name[:4]
		} else {
			name = name + strings.Repeat(" ", 4-len(name))
		}

		// Convert characters to bytes
		for charIdx := 0; charIdx < 4; charIdx++ {
			pageData[setupIdx*4+charIdx] = byte(name[charIdx])
		}
	}

	return []SysexPage{{
		Address: EEPROM_ADR_SETUP_NAMES,
		Data:    pageData,
	}}
}

// generateGroupNamesPages generates the GROUP_NAMES pages
func (sg *SysexGenerator) generateGroupNamesPages() []SysexPage {
	var pages []SysexPage
	address := uint16(EEPROM_ADR_GROUP_NAMES)

	// Each setup has 1 page for group names
	for setupIdx := 0; setupIdx < 16; setupIdx++ {
		pageData := make([]byte, 64)

		// Each group name is 4 characters
		for groupIdx := 0; groupIdx < 16; groupIdx++ {
			var name string
			if setupIdx < len(sg.config.Setups) && groupIdx < len(sg.config.Setups[setupIdx].Groups) {
				name = sg.config.Setups[setupIdx].Groups[groupIdx].Name
			}

			// Pad or truncate to 4 characters
			if len(name) > 4 {
				name = name[:4]
			} else {
				name = name + strings.Repeat(" ", 4-len(name))
			}

			// Convert characters to bytes
			for charIdx := 0; charIdx < 4; charIdx++ {
				pageData[groupIdx*4+charIdx] = byte(name[charIdx])
			}
		}

		pages = append(pages, SysexPage{
			Address: address,
			Data:    pageData,
		})
		address += 64
	}

	return pages
}

// generateSetupDataPages generates the SETUP_DATA pages
func (sg *SysexGenerator) generateSetupDataPages() []SysexPage {
	var pages []SysexPage
	address := uint16(EEPROM_ADR_SETUP_DATA)

	// Each setup has 48 pages (3072 bytes), each group uses 3 pages (192 bytes)
	for setupIdx := 0; setupIdx < 16; setupIdx++ {
		for pageInSetup := 0; pageInSetup < 48; pageInSetup++ {
			pageData := make([]byte, 64)

			// Each page contains data for 1/3 of a group (64 bytes)
			groupIdx := pageInSetup / 3
			groupPart := pageInSetup % 3

			if groupIdx < 16 {
				// Get the actual group data if available
				var group *Group
				if setupIdx < len(sg.config.Setups) && groupIdx < len(sg.config.Setups[setupIdx].Groups) {
					group = &sg.config.Setups[setupIdx].Groups[groupIdx]
				}

				// Generate the appropriate part of the group data
				groupData := sg.generateGroupDataPart(group, groupPart)
				copy(pageData, groupData)
			}

			pages = append(pages, SysexPage{
				Address: address,
				Data:    pageData,
			})
			address += 64
		}
	}

	return pages
}

// generateGroupDataPart generates a specific part of group data
func (sg *SysexGenerator) generateGroupDataPart(group *Group, part int) []byte {
	data := make([]byte, 64)

	if group == nil {
		return data
	}

	switch part {
	case 0: // First 64 bytes: encoder type/channel, link/command, command high, lower, upper
		// Encoder type/channel (bytes 0-15)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				encoderType := sg.convertEncoderType(setting.EC.Type)
				channel := setting.EC.Channel - 1 // Convert to 0-based
				if channel < 0 {
					channel = 0
				}
				data[i] = byte((encoderType << 4) | channel)
			}
		}

		// Encoder link + command (bytes 16-31)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				link := 0 // No link for now
				command := setting.EC.Number
				if command > 127 {
					command = 127
				}
				data[16+i] = byte((link << 7) | command)
			}
		}

		// Encoder command high (bytes 32-47)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				// For 14-bit CC, this would be the MSB
				data[32+i] = 0
			}
		}

		// Encoder lower values (bytes 48-63)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				data[48+i] = byte(setting.EC.Lower)
			}
		}

	case 1: // Second 64 bytes: upper, mode/scale, upper/lower msb, key type/channel
		// Encoder upper values (bytes 0-15)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				data[i] = byte(setting.EC.Upper)
			}
		}

		// Encoder mode/scale (bytes 16-31)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				mode := sg.convertEncoderMode(setting.EC.Mode)
				scale := sg.convertDisplayScale(setting.EC.Display)
				data[16+i] = byte((mode << 4) | scale)
			}
		}

		// Encoder upper/lower msb (bytes 32-47)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				// For 14-bit values, extract MSB
				upperMsb := (setting.EC.Upper >> 8) & 0x0F
				lowerMsb := (setting.EC.Lower >> 8) & 0x0F
				data[32+i] = byte((upperMsb << 4) | lowerMsb)
			}
		}

		// Push button type/channel (bytes 48-63)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				pbType := sg.convertPushButtonType(setting.PB.Type)
				channel := setting.PB.Channel - 1 // Convert to 0-based
				if channel < 0 {
					channel = 0
				}
				data[48+i] = byte((pbType << 4) | channel)
			}
		}

	case 2: // Third 64 bytes: encoder names
		// Encoder names (4 characters each, bytes 0-63)
		for i := 0; i < 16; i++ {
			if i < len(group.Settings) {
				setting := group.Settings[i]
				name := setting.Name

				// Pad or truncate to 4 characters
				if len(name) > 4 {
					name = name[:4]
				} else {
					name = name + strings.Repeat(" ", 4-len(name))
				}

				// Convert characters to bytes
				for charIdx := 0; charIdx < 4; charIdx++ {
					data[i*4+charIdx] = byte(name[charIdx])
				}
			}
		}
	}

	return data
}

// generateSetupKey2Pages generates the SETUP_KEY2 pages
func (sg *SysexGenerator) generateSetupKey2Pages() []SysexPage {
	var pages []SysexPage
	address := uint16(EEPROM_ADR_SETUP_KEY2)

	// Each setup has 8 pages (512 bytes), each group uses 0.5 pages (32 bytes)
	for setupIdx := 0; setupIdx < 16; setupIdx++ {
		for pageInSetup := 0; pageInSetup < 8; pageInSetup++ {
			pageData := make([]byte, 64)

			// Each page contains data for 2 groups
			for groupInPage := 0; groupInPage < 2; groupInPage++ {
				groupIdx := pageInSetup*2 + groupInPage
				if groupIdx >= 16 {
					break
				}

				// Get the actual group data if available
				var group *Group
				if setupIdx < len(sg.config.Setups) && groupIdx < len(sg.config.Setups[setupIdx].Groups) {
					group = &sg.config.Setups[setupIdx].Groups[groupIdx]
				}

				// Generate push button display + lower and link + upper data for this group
				groupData := sg.generatePushButtonKey2Data(group)

				// Copy group data to page (32 bytes per group)
				copy(pageData[groupInPage*32:], groupData)
			}

			pages = append(pages, SysexPage{
				Address: address,
				Data:    pageData,
			})
			address += 64
		}
	}

	return pages
}

// generatePushButtonKey2Data generates push button display + lower and link + upper data for a group
func (sg *SysexGenerator) generatePushButtonKey2Data(group *Group) []byte {
	data := make([]byte, 32)

	if group == nil {
		return data
	}

	// First 16 bytes: display + lower value
	for i := 0; i < 16; i++ {
		if i < len(group.Settings) {
			setting := group.Settings[i]

			// Convert display to bit
			display := 0
			if setting.PB.Display == "On" {
				display = 1
			}

			// Lower value (0-127)
			lower := setting.PB.Lower
			if lower > 127 {
				lower = 127
			}

			// Combine display and lower
			data[i] = byte((display << 7) | lower)
		}
	}

	// Second 16 bytes: link + upper value
	for i := 0; i < 16; i++ {
		if i < len(group.Settings) {
			setting := group.Settings[i]

			// Link (0 for now)
			link := 0

			// Upper value (0-127)
			upper := setting.PB.Upper
			if upper > 127 {
				upper = 127
			}

			// Combine link and upper
			data[16+i] = byte((link << 7) | upper)
		}
	}

	return data
}

// encodePage encodes a single page into sysex format
func (sg *SysexGenerator) encodePage(page SysexPage) []byte {
	var sysex []byte

	// Page address
	sysex = append(sysex, encodeCommand(CMD_PAGE_NUM_H, byte(page.Address>>8))...)
	sysex = append(sysex, encodeCommand(CMD_PAGE_NUM_L, byte(page.Address&0xFF))...)

	// Page data
	for _, b := range page.Data {
		sysex = append(sysex, encodeCommand(CMD_PAGE_DATA, b)...)
	}

	// Calculate CRC
	crc := sg.calculateCRC(page.Data)
	sysex = append(sysex, encodeCommand(CMD_PAGE_CRC_H, byte(crc>>8))...)
	sysex = append(sysex, encodeCommand(CMD_PAGE_CRC_L, byte(crc&0xFF))...)

	// Add 30 zero bytes
	for i := 0; i < 30; i++ {
		sysex = append(sysex, 0x00)
	}

	return sysex
}

// calculateCRC calculates the CRC for the data bytes
func (sg *SysexGenerator) calculateCRC(data []byte) uint16 {
	var crc uint16
	for _, b := range data {
		crc += uint16(b)
	}
	return crc
}

// encodeCommand encodes a command with data into sysex format
func encodeCommand(cmd byte, data byte) []byte {
	return []byte{
		cmd,                     // Command
		(2 << 4) | (data >> 4),  // Data MSB
		(1 << 4) | (data & 0xF), // Data LSB
	}
}

// Conversion functions
func (sg *SysexGenerator) convertEncoderType(typeStr string) int {
	switch typeStr {
	case "CCR1":
		return 0
	case "CCR2":
		return 1
	case "CCAb":
		return 2
	case "PrgC":
		return 3
	case "CCAh":
		return 4
	case "PBnd":
		return 5
	case "AftT":
		return 6
	case "Note":
		return 7
	case "NRPN":
		return 8
	default:
		return 2 // Default to CCAb
	}
}

func (sg *SysexGenerator) convertEncoderMode(modeStr string) int {
	switch modeStr {
	case "Div8":
		return 0
	case "Div4":
		return 1
	case "Div2":
		return 2
	case "Acc0":
		return 3
	case "Acc1":
		return 4
	case "Acc2":
		return 5
	case "Acc3":
		return 6
	case "LSp2":
		return 7
	case "LSp4":
		return 8
	case "LSp6":
		return 9
	default:
		return 3 // Default to Acc0
	}
}

func (sg *SysexGenerator) convertDisplayScale(display string) int {
	switch display {
	case "off":
		return 0
	case "127":
		return 1
	case "100":
		return 2
	case "1000":
		return 3
	case "±63":
		return 4
	case "±50":
		return 5
	case "±500":
		return 6
	case "ONOF":
		return 7
	case "9999":
		return 8
	default:
		return 1 // Default to 127
	}
}

func (sg *SysexGenerator) convertPushButtonType(typeStr string) int {
	switch typeStr {
	case "Off":
		return 0
	case "Note":
		return 1
	case "CC":
		return 2
	case "PrgC":
		return 3
	case "PBnd":
		return 4
	case "AftT":
		return 5
	case "Grp":
		return 6
	case "Set":
		return 7
	case "Acc3":
		return 8
	case "LSp6":
		return 9
	case "Min":
		return 10
	case "Max":
		return 11
	default:
		return 2 // Default to CC
	}
}

func (sg *SysexGenerator) convertPushButtonTypeToMode(modeStr string) int {
	switch modeStr {
	case "Key":
		return 0
	case "Togl":
		return 1
	default:
		return 0 // Default to Key
	}
}

// PrintSysexHex prints the sysex message in hex format
func PrintSysexHex(sysex []byte) {
	fmt.Println("Sysex message:")
	fmt.Println(hex.EncodeToString(sysex))
	fmt.Printf("Length: %d bytes\n", len(sysex))
}
