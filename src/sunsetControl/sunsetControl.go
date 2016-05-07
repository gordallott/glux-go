package sunsetControl

import (
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
	if currentTime.Before(sunsetTime) && currentTime.Before(duskTime) {
		SecondsUntilNextEvent = sunsetTime.Sub(currentTime).Seconds()
	} else if currentTime.After(sunsetTime) && currentTime.Before(duskTime) {
		SecondsUntilNextEvent = 60.0
	} else {
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

	return int(math.Ceil(newSunriseTime.Sub(currentTime).Seconds()));
}

func TimeOfDayBrightnessCalc(currentTime time.Time) float64 {
	sunTimes := suncalc.SunTimes(currentTime, AUTOLIGHT_LAT, AUTOLIGHT_LONG)
	sunsetTime := sunTimes["sunset"]
	duskTime := sunTimes["dusk"]

	if currentTime.Before(sunsetTime) {
		return 0.0
	}

	if currentTime.After(duskTime) {
		return 1.0
	}

	sunsetDuration := duskTime.Sub(sunsetTime).Seconds()
	secondsSinceSunset := currentTime.Sub(sunsetTime).Seconds()

	return secondsSinceSunset / sunsetDuration
}
