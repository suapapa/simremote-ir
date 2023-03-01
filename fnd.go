package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/suapapa/go_devices/tm1638"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const (
	blank = "        "
)

type FND struct {
	m          *tm1638.Module
	lastString string

	sync.Mutex
}

func NewFND() *FND {
	m, err := tm1638.Open(
		gpioreg.ByName("17"), // data
		gpioreg.ByName("27"), // clk
		gpioreg.ByName("22"), // stb
	)
	if err != nil {
		log.Fatal(err)
	}
	return &FND{
		m: m,
	}
}

func (f *FND) SetString(s string) {
	f.Lock()
	defer f.Unlock()
	if len(s) > 8 {
		s = s[:8]
	} else if len(s) < 8 {
		s += blank[:8-len(s)]
	}

	if f.lastString == s {
		return
	}

	f.m.SetString(s)
	f.lastString = s
}

func (f *FND) Welcome() error {
	f.Lock()
	defer f.Unlock()
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
loop:
	for i := 0; i < 100; i++ {
		f.m.SetString(randString(rnd))
		time.Sleep(30 * time.Millisecond)
	}
	_, ip, _, err := resolveNet()
	if err != nil || ip == "" {
		goto loop
	}
	f.m.SetString("srmt-" + ip[len(ip)-3:])
	return nil
}

func randString(rnd *rand.Rand) string {
	// gen 32~126
	var randString []byte
	for i := 0; i < 8; i++ {
		ch := byte(32 + rnd.Intn(126-32))
		randString = append(randString, ch)
	}
	return string(randString)
}

// resolveNet returns hostname, IP, MAC and error
func resolveNet() (string, string, string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", "", "", err
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", "", err
	}

	var ip net.IP
	for _, i := range ifaces {
		if (i.Flags&net.FlagUp) == 0 ||
			(i.Flags&net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// sometimes ip.To4() makes ip to nil
			if ip != nil {
				ip = ip.To4()
			}
			if ip != nil {
				return hostname, ip.String(), i.HardwareAddr.String(), nil
			}
		}
	}
	return "", "", "", fmt.Errorf("cannot resolve the IP")
}
