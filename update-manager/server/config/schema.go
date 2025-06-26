package config

type Schedule string

const (
	ASAP   Schedule = "asap"
	Daily  Schedule = "daily"
	Weekly Schedule = "weekly"
)

var ValidSchedules = []Schedule{
	ASAP,
	Daily,
	Weekly,
}

type UpdateType string

const (
	Major UpdateType = "major"
	Minor UpdateType = "minor"
	Patch UpdateType = "patch"
)

var ValidUpdateTypes = []UpdateType{
	Major,
	Minor,
	Patch,
}

type Group struct {
	Name         string       `yaml:"name"`
	Schedule     Schedule     `yaml:"schedule"`
	UpdateTypes  []UpdateType `yaml:"update_types"`
	SkipApproval bool         `yaml:"skip_approval"`
	Entities     []string     `yaml:"entities"`
}

type Catchall struct {
	Name         string   `yaml:"name"`
	Schedule     Schedule `yaml:"schedule"`
	SkipApproval bool     `yaml:"skip_approval"`
}

type Config struct {
	Groups []Group `yaml:"groups,omitempty"`

	Catchall *Catchall `yaml:"catchall,omitempty"`
}
