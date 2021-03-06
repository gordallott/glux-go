package main

import (
	"castcontrol"
	"fmt"
	"hueControl"
	"log"
	"math"
	"plexControl"
	"ps4control"
	"sunsetControl"
	"time"
)

const (
	BrightnessDim     = 0.2
	BrightnessNominal = 0.6
	BrightnessFull    = 0.8
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

			log.Printf("glux: targetBrightness: %v\n", targetBrightness)

			hueBridge := hueControl.GetHueBridge()
			lightsState, _ := hueControl.AreLightsOn(hueBridge)
			lightsBrightness, _ := hueControl.GetLightsBrightness(hueBridge)
			log.Printf("glux: actualbrightnes: %v\n", lightsBrightness)

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

		secondsUntilSunsetEvent := sunsetControl.SecondsUntilSunsetEvent(time.Now())
		secondsUntilSunriseEvent := sunsetControl.SecondsUntilSunriseEvent(time.Now())
		plexState := plexControl.GetPlexState()
		ps4State := ps4control.GetPs4State()
		castState := castcontrol.GetCastState()

		fmt.Printf("Lights state is %v\n", lightsState)
		fmt.Printf("seconds until sunset event: %v\n", secondsUntilSunsetEvent)
		fmt.Printf("seconds until sunrise event: %v\n", secondsUntilSunriseEvent)
		fmt.Printf("plex state: %v\n", plexState)
		fmt.Printf("ps4 state: %v\n", ps4State)
		fmt.Printf("cast state: %v\n", castState)

		// start of actual logic for lights
		// if lights are off, do nothing. single exception is if it is close to sunset event

		if secondsUntilSunsetEvent < 120 && secondsUntilSunsetEvent >= 0 && lightsState == false {
			log.Printf("only %v seconds until sunset event, turning lights on", secondsUntilSunsetEvent)
			hueControl.TurnLightsOn(hueBridge)
			lightsState, _ = hueControl.AreLightsOn(hueBridge)
		}

		// if the lights are on, we probably want to turn them off at sunrise, we don't do this in a clever way, they just turn off instantly
		if secondsUntilSunriseEvent < 120 && secondsUntilSunriseEvent >= 0 {
			log.Printf("only %v seconds until sunrise event, turning lights off")
			hueControl.TurnLightsOff(hueBridge)
			lightsState, _ = hueControl.AreLightsOn(hueBridge)
		}

		if lightsState == false {
			continue
		}

		// if we get here then the lights are on and we should start doing clever things

		timeOfDayBrightness := sunsetControl.TimeOfDayBrightnessCalc(time.Now())
		var plexBrightness float64
		switch plexControl.GetPlexState() {
		case plexControl.StatePlaying:
			plexBrightness = 0.0
		case plexControl.StatePaused:
			plexBrightness = 0.5
		case plexControl.StateStopped:
			plexBrightness = 1.0
		}

		var ps4Brightness float64
		switch ps4control.GetPs4State() {
		case ps4control.StatePlaying:
			ps4Brightness = 0.0
		case ps4control.StatePaused:
			ps4Brightness = 0.5
		case ps4control.StateStopped:
			ps4Brightness = 1.0
		}

		var castBrightness float64
		switch castcontrol.GetCastState() {
		case castcontrol.StatePlaying:
			castBrightness = 0.0
		case castcontrol.StatePaused:
			castBrightness = 0.5
		case castcontrol.StateStopped:
			castBrightness = 1.0
		}

		fmt.Printf("timeOfDayBrightness: %v\n", timeOfDayBrightness)
		fmt.Printf("plexBrightness: %v\n", plexBrightness)
		fmt.Printf("ps4Brightness: %v\n", ps4Brightness)
		fmt.Printf("castBrightness: %v\n", castBrightness)

		var combBrightness = math.Min(plexBrightness, ps4Brightness)
		combBrightness = math.Min(combBrightness, castBrightness)

		totalBrightness := timeOfDayBrightness * combBrightness
		brightnessMessages <- totalBrightness * BrightnessFull

		fmt.Printf("total brightness: %v\n", totalBrightness*BrightnessFull)
	}
}

func main() {
	fmt.Printf("Starting glux-go...\n")
	mainLoopFunc()
}
