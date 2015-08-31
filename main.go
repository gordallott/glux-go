package main

import (
	"fmt"
	"hueControl"
	"log"
	"math"
	"plexControl"
	"sunsetControl"
	"time"
)

const (
	BrightnessDim     = 0.2
	BrightnessNominal = 0.4
	BrightnessFull    = 0.6
	BrightnessTick    = 0.02
)

func closeTo(a float64, b float64) bool {
	return math.Abs(a-b) < BrightnessTick*1.5
}

func animationLoop(brightnessValues chan float64) {
	// loops forever, if there are no new brightness values coming in and no animations to do then it
	// blocks on brightnessValues waiting for data.
	var targetBrightness float64
	for ; ; targetBrightness = <-brightnessValues { // outer loop blocks
		for ; ; time.Sleep(250 * time.Millisecond) { // inner loop does not block
			select {
			case targetBrightness = <-brightnessValues:
			default:
			}

			hueBridge := hueControl.GetHueBridge()
			lightsState, _ := hueControl.AreLightsOn(hueBridge)
			lightsBrightness, _ := hueControl.GetLightsBrightness(hueBridge)

			if lightsState == true && closeTo(lightsBrightness, targetBrightness) == false {

				change := BrightnessTick
				if targetBrightness < lightsBrightness {
					change = change - change - change // inverts positive to negative
				}

				hueControl.SetBrightness(hueBridge, lightsBrightness+change)
			} else {
				// either lights are off or we reached our target, either way stop doing things
				// and bubble out to the blocking loop
				break
			}

		}
	}
}

func mainLoopFunc() {
	brightnessMessages := make(chan float64)
	go animationLoop(brightnessMessages)
	for ; ; time.Sleep(5 * time.Second) {
		hueBridge := hueControl.GetHueBridge()
		lightsState, _ := hueControl.AreLightsOn(hueBridge)

		secondsUntilSunsetEvent := sunsetControl.SecondsUntilNextEvent(time.Now())
		plexState := plexControl.GetPlexState()

		fmt.Printf("Lights state is %v\n", lightsState)
		fmt.Printf("seconds until sunset event: %v\n", secondsUntilSunsetEvent)
		fmt.Printf("plex state: %v\n", plexState)

		// start of actual logic for lights
		// if lights are off, do nothing. single exception is if it is close to sunset event

		if secondsUntilSunsetEvent < 120 && secondsUntilSunsetEvent >= 0 && lightsState == false {
			log.Printf("only %v seconds until sunset event, turning lights on", secondsUntilSunsetEvent)
			hueControl.TurnLightsOn(hueBridge)
			lightsState, _ = hueControl.AreLightsOn(hueBridge)
		}

		if lightsState == false {
			continue
		}

		// if we get here then the lights are on and we should start doing clever things

		timeOfDayBrightness := sunsetControl.TimeOfDayBrightnessCalc(time.Now()) * BrightnessFull
		var plexBrightness float64
		switch plexControl.GetPlexState() {
		case plexControl.StatePlaying:
			plexBrightness = BrightnessDim
		case plexControl.StatePaused:
			plexBrightness = BrightnessNominal
		case plexControl.StateStopped:
			plexBrightness = BrightnessFull
		}

		totalBrightness := timeOfDayBrightness * plexBrightness
		brightnessMessages <- totalBrightness
	}
}

func main() {
	fmt.Printf("Starting glux-go...\n")
	mainLoopFunc()
}
