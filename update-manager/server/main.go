package main

import (
	"fmt"

	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
	"github.com/lucaspopp0/hass-update-manager/update-manager/util"
)

func main() {
	hass := homeassistant.NewAPI(homeassistant.APIConfig{
		SupervisorToken: util.GetEnv("SUPERVISOR_TOKEN", ""),
	})

	states, err := hass.GetStates()
	if err != nil {
		panic(err)
	}

	fmt.Println(util.MarshalIndent(states))
}
