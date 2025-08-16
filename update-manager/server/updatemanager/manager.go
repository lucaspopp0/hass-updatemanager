package updatemanager

import (
	"fmt"
	"time"

	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
)

type Config struct {
	HomeAssistant homeassistant.API
}

type Manager interface {
	Run() error
	CheckForUpdates() ([]homeassistant.UpdateEntity, error)
}

type manager struct {
	Config
}

func NewManager(cfg Config) Manager {
	return &manager{cfg}
}

func (m *manager) Run() error {
	for {
		updates, err := m.CheckForUpdates()
		if err != nil {
			return err
		}

		fmt.Printf("%v updates available\n", len(updates))

		time.Sleep(time.Minute)
	}
}

func (m *manager) CheckForUpdates() ([]homeassistant.UpdateEntity, error) {
	updates, err := m.HomeAssistant.ListUpdates()
	if err != nil {
		return nil, err
	}

	canUpdate := []homeassistant.UpdateEntity{}

	for _, update := range updates {
		if update.UpdateAvailable() {
			canUpdate = append(canUpdate, update)
		}
	}

	return canUpdate, nil
}
