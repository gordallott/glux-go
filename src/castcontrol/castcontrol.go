package castcontrol

import (
	"fmt"
	"github.com/hashicorp/mdns"
	"strings"
	"time"
)

const (
	castService    = "_googlecast._tcp"
	castType       = "Chromecast"
	castDeviceName = "Living Room"
)

const (
	StatePlaying = iota
	StatePaused
	StateStopped
)

func GetCastState() int {
	// Make a channel for results and start listening
	entries := make(chan *mdns.ServiceEntry, 4)

	go func() {
		mdns.Query(&mdns.QueryParam{
			Service: castService,
			Timeout: time.Second * 3,
			Entries: entries,
		})
		close(entries)
	}()

	for entry := range entries {
		for _, field := range entry.InfoFields {
			if strings.HasPrefix(field, "md=") && field != fmt.Sprintf("md=%s", castType) {
				fmt.Printf("Ignoring %s (wrong device type)\n", entry.Name)
				continue
			}

			if strings.HasPrefix(field, "fn=") && field != fmt.Sprintf("fn=%s", castDeviceName) {
				fmt.Printf("Ignoring %s (wrong device name)\n", entry.Name)
				continue
			}

			if strings.Contains(field, "st=1") {
				return StatePlaying
			}
		}
	}

	return StateStopped
}
