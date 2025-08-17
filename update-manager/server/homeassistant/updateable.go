package homeassistant

import (
	"encoding/json"
	"fmt"
	"strings"
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
	return u.State == "on"
}

func isUpdate(entityState EntityState[map[string]any]) (*UpdateEntity, bool) {
	if !strings.HasPrefix(entityState.EntityID, "update.") {
		return nil, false
	}

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

func (c *apiClient) GetUpdate(entityID string) (*UpdateEntity, error) {
	entity, err := c.GetState(entityID)
	if err != nil {
		return nil, err
	}

	update, ok := isUpdate(*entity)
	if !ok {
		return nil, fmt.Errorf("%v is not an update", entity.EntityID)
	}

	return update, nil
}

func (c *apiClient) InstallUpdates(entityIDs []string) error {
	_, err := c.CallService("update/install", map[string]any{
		"entity_id": entityIDs,
	})

	return err
}

func (c *apiClient) Restart() error {
	_, err := c.CallService("homeassistant/restart", map[string]any{})

	return err
}
