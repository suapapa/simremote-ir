package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/suapapa/go_devices/tm1638"
)

var (
	rnd *rand.Rand
)

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func displayWelcome(dev *tm1638.Module) error {
loop:
	for i := 0; i < 100; i++ {
		dev.SetString(randString())
		time.Sleep(30 * time.Millisecond)
	}
	_, ip, _, err := resolveNet()
	if err != nil || ip == "" {
		goto loop
	}
	dev.SetString("srmt-" + ip[len(ip)-3:])
	return nil
}

func randString() string {
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
