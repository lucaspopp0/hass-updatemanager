package updatemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
	"golang.org/x/sync/errgroup"
)

type Manager interface {
	Run() error
	CheckForUpdates() ([]homeassistant.UpdateEntity, error)
}

type Config struct {
	HomeAssistant homeassistant.API

	MaintenanceDetails MaintenanceDetails
}

type MaintenanceDetails struct {
	StartTime time.Duration
	Duration  time.Duration
}

type manager struct {
	Config

	// A tracker for the number of available updates
	availableUpdates int
}

func NewManager(cfg Config) Manager {
	return &manager{
		Config: cfg,

		availableUpdates: -1,
	}
}

func (m *manager) Run() error {
	for {
		updates, err := m.CheckForUpdates()
		if err != nil {
			return err
		}

		availableUpdates := len(updates)
		if availableUpdates != m.availableUpdates {
			err = m.HomeAssistant.PostState(
				homeassistant.EntityState[map[string]any]{
					EntityID: "sensor.updatemanager_available_updates",
					State:    fmt.Sprintf("%d", len(updates)),
					Attributes: map[string]any{
						"unit_of_measurement": "updates",
						"friendly_name":       "Available Updates",
					},
				},
			)

			if err != nil {
				return err
			}
		}

		if m.canMaintenance() && len(updates) > 0 {
			fmt.Printf("%v updates available. Installing now...\n", len(updates))

			err = m.installUpdates(updates)
			if err != nil {
				return err
			}

			fmt.Println("Restarting...")
			return m.HomeAssistant.Restart()
		}

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

func (m *manager) canMaintenance() bool {
	now := time.Now()
	nowTime := time.Hour*time.Duration(now.Hour()) +
		time.Minute*time.Duration(now.Minute())

	return m.MaintenanceDetails.StartTime <= nowTime &&
		nowTime <= m.MaintenanceDetails.StartTime+m.MaintenanceDetails.Duration
}

func (m *manager) installUpdates(updates []homeassistant.UpdateEntity) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	group, _ := errgroup.WithContext(ctx)

	entityIDs := make([]string, len(updates))
	for i, update := range updates {
		entityIDs[i] = update.EntityID

		group.Go(func() error {
			return m.waitForUpdate(ctx, update.EntityID)
		})
	}

	err := m.HomeAssistant.InstallUpdates(entityIDs)
	if err != nil {
		return err
	}

	err = group.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) waitForUpdate(ctx context.Context, entityID string) error {
	err := m.waitForUpdateStart(ctx, entityID)
	if err != nil {
		return err
	}

	fmt.Printf("update %q started\n", entityID)

	err = m.waitForUpdateFinish(ctx, entityID)
	if err != nil {
		return err
	}

	fmt.Printf("update %q finished\n", entityID)
	return nil
}

func (m *manager) waitForUpdateStart(ctx context.Context, entityID string) error {
	// Poll every 200ms for the update to begin
	ticker := time.NewTicker(200 * time.Millisecond)

	// Timeout if the update has not begun after 15s
	timeout := time.NewTimer(15 * time.Second)
	defer timeout.Stop()

	started := false
	for !started {
		select {

		case <-ctx.Done():
			return context.Canceled

		case <-timeout.C:
			return fmt.Errorf("timed out waiting for update to begin")

		case <-ticker.C:

		}

		update, err := m.HomeAssistant.GetUpdate(entityID)
		if err != nil {
			return err
		}

		if update.Attributes.InProgress {
			started = true
		}
	}

	return nil
}

func (m *manager) waitForUpdateFinish(ctx context.Context, entityID string) error {
	// Poll every second for the update to complete
	ticker := time.NewTicker(time.Second)

	// Timeout if the update has not finished after 5 minutes
	timeout := time.NewTimer(5 * time.Minute)
	defer timeout.Stop()

	finished := false
	for !finished {
		select {

		case <-ctx.Done():
			return context.Canceled

		case <-timeout.C:
			return fmt.Errorf("timed out waiting for update to begin")

		case <-ticker.C:

		}

		update, err := m.HomeAssistant.GetUpdate(entityID)
		if err != nil {
			return err
		}

		if !update.Attributes.InProgress {
			finished = true
		}
	}

	return nil
}
