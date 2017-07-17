package sunsetControl

import (
	"fmt"
	"github.com/muja/suncalc-go"
	"math"
	"time"
)

const (
	AUTOLIGHT_LAT  = 53.100405
	AUTOLIGHT_LONG = -2.443821
)

func SecondsUntilSunsetEvent(currentTime time.Time) int {
	sunTimes := suncalc.SunTimes(currentTime, AUTOLIGHT_LAT, AUTOLIGHT_LONG)
	sunsetTime := sunTimes["sunset"]
	duskTime := sunTimes["dusk"]

	var SecondsUntilNextEvent float64
	if currentTime.Before(sunsetTime) && currentTime.Before(duskTime) { // before sunset starts
		SecondsUntilNextEvent = sunsetTime.Sub(currentTime).Seconds()
	} else if currentTime.After(sunsetTime) && currentTime.Before(duskTime) { // during sunset
		SecondsUntilNextEvent = -1
	} else { // after sunset
		/*var tomorrow = time.Date(currentTime.Year(),
			currentTime.Month(),
			currentTime.Day()+1,
			0, 0, 0, 0, currentTime.Location())
		SecondsUntilNextEvent = tomorrow.Sub(currentTime).Seconds() + 60.0*/
		SecondsUntilNextEvent = -1
	}

	return int(math.Ceil(SecondsUntilNextEvent))
}

func SecondsUntilSunriseEvent(currentTime time.Time) int {
	sunTimes := suncalc.SunTimes(currentTime, AUTOLIGHT_LAT, AUTOLIGHT_LONG)
	sunriseTime := sunTimes["dawn"]

	year, month, day := sunriseTime.Date()
	newSunriseTime := time.Date(year, month, day, 2, 30, 0, 0, currentTime.Location())

	return int(math.Ceil(newSunriseTime.Sub(currentTime).Seconds()))
}

func TimeOfDayBrightnessCalc(currentTime time.Time) float64 {
	sunTimes := suncalc.SunTimes(currentTime, AUTOLIGHT_LAT, AUTOLIGHT_LONG)
	sunsetTime := sunTimes["sunset"]
	sunriseTime := sunTimes["sunrise"]
	lightsOnStartTime := sunsetTime.Add(-time.Hour)
	duskTime := sunTimes["dusk"]

	fmt.Printf("sunsettime: %s\n", sunsetTime)
	fmt.Printf("duskTime: %s\n", duskTime)
	fmt.Printf("sunriseTime: %s\n", sunriseTime)

	if currentTime.Before(lightsOnStartTime) && currentTime.After(sunriseTime) {
		return 0.0
	}

	if currentTime.After(duskTime) || currentTime.Before(sunriseTime) {
		return 1.0
	}

	sunsetDuration := duskTime.Sub(lightsOnStartTime).Seconds()
	secondsSinceSunset := currentTime.Sub(lightsOnStartTime).Seconds()

	fmt.Printf("SunsetDuration: %v\n SecondsSinceSunset: %v\n", sunsetDuration, secondsSinceSunset)

	return secondsSinceSunset / sunsetDuration
}
