package hueControl

import (
	"github.com/gordallott/hue"
	"log"
	"math"
)

const (
	HUE_IP       = "192.168.1.27"
	HUE_USERNAME = "d8cc566111b90642a0bffb238b703df"
)

func AreLightsOn(hueBridge *hue.Hue) (bool, error) {
	var resp hue.GetLightsResponse
	err := hueBridge.GetLights(&resp)
	if err != nil {
		return false, err
	}

	for light, _ := range resp {
		var lightResp hue.GetLightResponse
		err := hueBridge.GetLight(light, &lightResp)
		if err != nil {
			continue
		}
		if lightResp.State.On {
			return true, nil
		}
	}

	return false, nil
}

func GetLightsBrightness(hueBridge *hue.Hue) (float64, error) {
	var resp hue.GetLightsResponse
	err := hueBridge.GetLights(&resp)
	if err != nil {
		return 0.0, err
	}

	for light, _ := range resp {
		var lightResp hue.GetLightResponse
		err := hueBridge.GetLight(light, &lightResp)
		if err != nil {
			continue
		}
		return float64(lightResp.State.Bri) / 100.0, nil
	}

	return 0.0, nil
}

func TurnLightsOn(hueBridge *hue.Hue) error {
	return setBrightnessInternal(hueBridge, StateOn, 0.0)
}

func SetBrightness(hueBridge *hue.Hue, brightness float64) error {
	log.Printf("glux: setting brightness to %v\n", brightness)
	return setBrightnessInternal(hueBridge, StateOn, brightness)
}

const (
	StateOn = iota
	StateOff
)

func setBrightnessInternal(hueBridge *hue.Hue, onOffState int, brightness float64) error {
	var resp hue.GetLightsResponse
	err := hueBridge.GetLights(&resp)
	if err != nil {
		return err
	}

	var lightRequest hue.PutLightRequest
	stateOn := true
	stateBri := int(math.Max(math.Min(brightness*255, 255), 0))
	lightRequest.On = &stateOn
	lightRequest.Bri = &stateBri

	for light, _ := range resp {
		hueBridge.PutLight(light, &lightRequest)
	}

	return nil
}

func GetHueBridge() *hue.Hue {
	return &hue.Hue{HUE_IP, HUE_USERNAME, "glux-go"}
}
