package plexControl

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

const (
	PLEX_HOSTNAME = "http://192.168.1.14:32400"
)

const (
	StatePlaying = iota
	StatePaused
	StateStopped
)

func GetPlexState() int {

	re := regexp.MustCompile("state=\"(?P<state>(playing|paused|stopped))\" ")

	resp, err := http.Get(PLEX_HOSTNAME + "/status/sessions")
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error getting plex status: %s", err)
		return StateStopped
	}

	body, _ := ioutil.ReadAll(resp.Body)

	match := re.FindStringSubmatch(string(body))
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		result[name] = match[i]
	}

	if _, present := result["state"]; present == true {
		switch result["state"] {
		case "playing":
			return StatePlaying
		case "paused":
			return StatePaused
		case "stopped":
			return StateStopped
		}
	}

	return StateStopped

}
