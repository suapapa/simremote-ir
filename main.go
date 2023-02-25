package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/tarm/serial"
)

var (
	port       string
	skipRepeat bool
	apiAddr    string
)

func main() {
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "serial port")
	flag.BoolVar(&skipRepeat, "skip-repeat", false, "skip repeated codes")
	flag.StringVar(&apiAddr, "api", "http://localhost:5000", "api address")
	flag.Parse()

	ser, err := serial.OpenPort(&serial.Config{
		Name:        port,
		Baud:        9600,
		ReadTimeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer ser.Close()

	var code, lastCode uint32
	for {
		scanner := bufio.NewScanner(ser)
		for scanner.Scan() {
			codeStr := scanner.Text()
			if codeStr == "0" {
				if skipRepeat {
					continue
				}
				code = lastCode
			} else {
				codeBytes, err := hex.DecodeString(codeStr)
				if err != nil {
					log.Printf("error decoding %q: %v", codeStr, err)
					continue
				}
				code = uint32(codeBytes[0])<<24 | uint32(codeBytes[1])<<16 | uint32(codeBytes[2])<<8 | uint32(codeBytes[3])
			}

			log.Printf("code: 0x%x", code)
			lastCode = code

			button := modes[currMode][code]
			log.Printf("button: %s", button)

			if button == MODE {
				currMode = (currMode + 1) % len(modes)
				log.Printf("mode: %d", currMode)
			} else {
				if err := handle(button); err != nil {
					log.Printf("error doing %s: %v", button, err)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error scanning: %v", err)
		}
	}
}
