package main

import (
	"fmt"
	"time"

	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
	"github.com/lucaspopp0/hass-update-manager/update-manager/util"
)

func main() {
	hass := homeassistant.NewAPI(homeassistant.APIConfig{
		SupervisorToken: util.GetEnv("SUPERVISOR_TOKEN", ""),
	})

	updates, err := hass.ListUpdates()
	if err != nil {
		panic(err)
	}

	fmt.Println(util.MarshalIndent(updates))

	fmt.Printf("Current time: %v\n", time.Now().Format(time.DateTime))
}
