package main

import (
	"fmt"

	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
	"github.com/lucaspopp0/hass-update-manager/update-manager/updatemanager"
	"github.com/lucaspopp0/hass-update-manager/update-manager/util"
)

func main() {
	hass := homeassistant.NewAPI(homeassistant.APIConfig{
		SupervisorToken: util.GetEnv("SUPERVISOR_TOKEN", ""),
	})

	manager := updatemanager.NewManager(updatemanager.Config{
		HomeAssistant: hass,
	})

	err := manager.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
