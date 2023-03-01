package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
	"periph.io/x/host/v3"
)

var (
	port            string
	enableLongPress bool
	apiAddr         string
)

func main() {
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "serial port")
	flag.BoolVar(&enableLongPress, "l", false, "enable long press")
	flag.StringVar(&apiAddr, "api", "http://localhost:5000", "api address")
	flag.Parse()

	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	fnd := NewFND()
	fnd.Welcome()

	apiC := NewAPIClient(apiAddr)

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
		fnd.SetString(fmt.Sprintf("M%d", curMode))
		scanner := bufio.NewScanner(ser)
		for scanner.Scan() {
			codeStr := scanner.Text()
			if codeStr == "0" {
				if !enableLongPress {
					continue
				}
				code = lastCode
			} else {
				codeBytes, err := hex.DecodeString(codeStr)
				if err != nil {
					log.Printf("error decoding %q: %v", codeStr, err)
					continue
				}
				if len(codeBytes) != 4 {
					log.Printf("unexpected code length: %d", len(codeBytes))
					continue
				}
				code = uint32(codeBytes[0])<<24 | uint32(codeBytes[1])<<16 | uint32(codeBytes[2])<<8 | uint32(codeBytes[3])
			}

			log.Printf("code: 0x%x", code)
			lastCode = code

			button, ok := modes[curMode][code]
			if !ok {
				log.Printf("unknown code: 0x%x", code)
				continue
			}

			log.Printf("button: %s", button)
			fnd.SetString(fmt.Sprintf("M%d-%s", curMode, button))

			if button == MODE {
				curMode = (curMode + 1) % len(modes)
				log.Printf("mode: %d", curMode)
				fnd.SetString(fmt.Sprintf("M%d", curMode))
			} else {
				if err := apiC.Handle(button); err != nil {
					log.Printf("error doing %s: %v", button, err)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error scanning: %v", err)
		}
	}
}
