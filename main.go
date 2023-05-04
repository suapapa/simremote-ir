package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"log"
	"time"

	"github.com/tarm/serial"
	"periph.io/x/host/v3"
)

var (
	port    string
	apiAddr string
	dryRun  bool
)

func main() {
	flag.StringVar(&port, "port", "/dev/ttyUSB0", "serial port")
	flag.BoolVar(&dryRun, "n", false, "dry run")
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

	codeInC := make(chan uint32)
	buttonC := make(chan button)

	go func(codeC chan uint32) {
		var repeatedCode uint32
		var repeatCnt int
		var downTime time.Time
		timeTk := time.NewTicker(100 * time.Millisecond)
		defer timeTk.Stop()
		for {
			select {
			case code := <-codeC:
				if repeatedCode == code {
					repeatCnt++
				} else {
					repeatedCode = code
					repeatCnt = 1
				}
				downTime = time.Now()
			case <-timeTk.C:
				// 3회 이상 반복 시그널이 들어오고, 마지막으로 눌린 뒤 500ms가 지났으면
				// 키가 떼진 것으로 간주한다.
				if repeatedCode > 0 {
					if time.Since(downTime) > 500*time.Millisecond && repeatCnt >= 3 {
						log.Printf("code, 0x%x longPress", repeatedCode)
						btn := codeToButton(repeatedCode, true)
						buttonC <- btn
						repeatedCode = 0
						repeatCnt = 0
					} else if time.Since(downTime) > 100*time.Millisecond && repeatCnt < 3 {
						log.Printf("code, 0x%x shortPress", repeatedCode)
						btn := codeToButton(repeatedCode, false)
						buttonC <- btn
						repeatedCode = 0
						repeatCnt = 0
					}
				}
			}
		}
	}(codeInC)

	go func(buttonC chan button) {
		for btn := range buttonC {
			log.Printf("button: %s", btn)
			go func() {
				fnd.SetString(btn.String())
				time.Sleep(1 * time.Second)
				fnd.SetString("        ")
			}()
			if err := apiC.Handle(btn); err != nil {
				log.Printf("error doing %s: %v", btn, err)
			}
		}
	}(buttonC)

	// ir key scanner
	for {
		scanner := bufio.NewScanner(ser)
		for scanner.Scan() {
			codeStr := scanner.Text()
			if codeStr == "0" {
				codeInC <- lastCode
				continue
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
				lastCode = code
				codeInC <- code
			}
			if err := scanner.Err(); err != nil {
				log.Printf("error scanning: %v", err)
			}
		}
	}
}

func codeToButton(code uint32, longPress bool) button {
	log.Printf("code: 0x%x", code)
	// lastCode = code
	var codeMap map[uint32]button
	if !longPress {
		codeMap = shortPressButton
	} else {
		codeMap = longPressButton
	}

	btn, ok := codeMap[code]
	if !ok {
		log.Printf("unknown code: 0x%x", code)
		return UNKNOWN
	}
	return btn
}
