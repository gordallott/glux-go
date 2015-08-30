package main

import (
	"fmt"
	"hueControl"
	"plexControl"
	"sunsetControl"
	"time"
)

func main() {
	fmt.Printf("Starting glux-go...\n")

	hueBridge := hueControl.GetHueBridge()
	lightsState, _ := hueControl.AreLightsOn(hueBridge)

	secondsUntilSunsetEvent := sunsetControl.SecondsUntilNextEvent(time.Now())

	fmt.Printf("Lights state is %v\n", lightsState)
	fmt.Printf("seconds until sunset event: %v\n", secondsUntilSunsetEvent)

	plexState := plexControl.GetPlexState()
	fmt.Printf("plex state: %v", plexState)
}
