package homeassistant

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

type Update struct {
	EntityID     string
	FriendlyName string `yaml:"friendly_name"`

	AutoUpdate bool `yaml:"auto_update"`

	InstalledVersion string `yaml:"installed_version"`
	LatestVersion    string `yaml:"latest_version"`

	InProgress bool `yaml:"in_progress"`

	// More fields available
}

func isUpdate(entityState EntityState) (*Update, bool) {
	update := &Update{
		EntityID: entityState.EntityID,
	}

	attrBytes, err := yaml.Marshal(entityState.Attributes)
	if err != nil {
		fmt.Printf("error in marshaling: %v\n", err.Error())
		return nil, false
	}

	err = yaml.Unmarshal(attrBytes, &update)
	if err != nil {
		fmt.Printf("error in unmarshaling: %v\n", err.Error())
		return nil, false
	}

	if !strings.HasPrefix("update.", entityState.EntityID) {
		fmt.Printf("entity %q missing prefix\n", entityState.EntityID)
		return nil, false
	}

	return update, true
}

func (c *apiClient) ListUpdates() ([]Update, error) {
	entityStates, err := c.GetStates()
	if err != nil {
		return nil, err
	}

	updates := []Update{}
	for _, es := range entityStates {
		if update, ok := isUpdate(es); ok {
			updates = append(updates, *update)
		}
	}

	return updates, nil
}

func (c *apiClient) InstallUpdate(entityID string) error {
	_, err := c.CallService("update/install", map[string]any{
		"entity_id": entityID,
	})

	return err
}
