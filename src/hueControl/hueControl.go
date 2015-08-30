package hueControl

import (
	"github.com/bklimt/hue"
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

	for key, _ := range resp {
		var lightResp hue.GetLightResponse
		err := hueBridge.GetLight(key, &lightResp)
		if err != nil {
			continue
		}
		if lightResp.State.On {
			return true, nil
		}
	}

	return false, nil
}

func GetHueBridge() *hue.Hue {
	return &hue.Hue{HUE_IP, HUE_USERNAME, "glux-go"}
}
