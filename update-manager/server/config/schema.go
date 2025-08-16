package config

import "github.com/lucaspopp0/hass-update-manager/update-manager/schedule"

type UpdateType string

const (
	Major UpdateType = "major"
	Minor UpdateType = "minor"
	Patch UpdateType = "patch"
)

type Group struct {
	Name         string            `yaml:"name"`
	Schedule     schedule.Schedule `yaml:"schedule"`
	SkipApproval bool              `yaml:"skip_approval"`
	Entities     []string          `yaml:"entities"`
}

type Catchall struct {
	Name         string            `yaml:"name"`
	Schedule     schedule.Schedule `yaml:"schedule"`
	SkipApproval bool              `yaml:"skip_approval"`
}

type Config struct {
	Groups []Group `yaml:"groups,omitempty"`

	Catchall *Catchall `yaml:"catchall,omitempty"`
}
