package homeassistant

import (
	"encoding/json"
	"fmt"
)

type UpdateEntity EntityState[UpdateAttributes]

type UpdateAttributes struct {
	FriendlyName string `json:"friendly_name"`

	AutoUpdate bool `json:"auto_update"`

	InstalledVersion string `json:"installed_version"`
	LatestVersion    string `json:"latest_version"`

	InProgress bool `json:"in_progress"`

	// More fields available
}

func (u UpdateEntity) UpdateAvailable() bool {
	return u.Attributes.InstalledVersion != u.Attributes.LatestVersion
}

func isUpdate(entityState EntityState[map[string]any]) (*UpdateEntity, bool) {
	entityBytes, err := json.Marshal(entityState)
	if err != nil {
		fmt.Printf("error in marshaling: %v\n", err.Error())
		return nil, false
	}

	update := UpdateEntity{}
	err = json.Unmarshal(entityBytes, &update)
	if err != nil {
		fmt.Printf("error in unmarshaling: %v\n", err.Error())
		return nil, false
	}

	return &update, true
}

func (c *apiClient) ListUpdates() ([]UpdateEntity, error) {
	entityStates, err := c.GetStates()
	if err != nil {
		return nil, err
	}

	updates := []UpdateEntity{}
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
