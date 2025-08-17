package main

import (
	"fmt"
	"time"

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

		MaintenanceDetails: updatemanager.MaintenanceDetails{
			StartTime: time.Duration(12+5) * time.Hour,
			Duration:  2 * time.Hour,
		},
	})

	err := manager.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
