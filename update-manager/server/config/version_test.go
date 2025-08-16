package config_test

import (
	"testing"

	"github.com/lucaspopp0/hass-update-manager/update-manager/config"
)

func TestVersionSegment_Compare(t *testing.T) {
	tests := []struct {
		name     string
		a        config.VersionSegment
		b        config.VersionSegment
		expected int
	}{
		{"equal strings", "1", "1", 0},
		{"equal numeric", "10", "10", 0},
		{"numeric less than", "1", "2", -1},
		{"numeric greater than", "10", "5", 1},
		{"numeric vs string", "1", "a", -1},
		{"string vs numeric", "a", "1", 1},
		{"string less than", "a", "b", -1},
		{"string greater than", "z", "a", 1},
		{"multi-digit numeric", "10", "2", 1},
		{"empty vs non-empty", "", "1", -1},
		{"non-empty vs empty", "1", "", 1},
		{"empty vs empty", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			if result != tt.expected {
				t.Errorf("VersionSegment.Compare(%q, %q) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected config.VersionString
	}{
		{
			"simple version",
			"1.2.3",
			config.VersionString{
				Prefix: "",
				Major:  "1",
				Minor:  "2",
				Patch:  "3",
				Suffix: "",
			},
		},
		{
			"version with prefix",
			"v1.2.3",
			config.VersionString{
				Prefix: "v",
				Major:  "1",
				Minor:  "2",
				Patch:  "3",
				Suffix: "",
			},
		},
		{
			"version with suffix",
			"1.2.3-beta",
			config.VersionString{
				Prefix: "",
				Major:  "1",
				Minor:  "2",
				Patch:  "3",
				Suffix: "-beta",
			},
		},
		{
			"version with prefix and suffix",
			"v1.2.3-rc1",
			config.VersionString{
				Prefix: "v",
				Major:  "1",
				Minor:  "2",
				Patch:  "3",
				Suffix: "-rc1",
			},
		},
		{
			"major only",
			"5",
			config.VersionString{
				Prefix: "",
				Major:  "5",
				Minor:  "",
				Patch:  "",
				Suffix: "",
			},
		},
		{
			"major.minor only",
			"1.5",
			config.VersionString{
				Prefix: "",
				Major:  "1",
				Minor:  "5",
				Patch:  "",
				Suffix: "",
			},
		},
		{
			"complex prefix",
			"release-1.2.3",
			config.VersionString{
				Prefix: "release-",
				Major:  "1",
				Minor:  "2",
				Patch:  "3",
				Suffix: "",
			},
		},
		{
			"no dots",
			"123",
			config.VersionString{
				Prefix: "",
				Major:  "123",
				Minor:  "",
				Patch:  "",
				Suffix: "",
			},
		},
		{
			"non-numeric segments",
			"1.2a.3b",
			config.VersionString{
				Prefix: "",
				Major:  "1",
				Minor:  "2a",
				Patch:  "3b",
				Suffix: "",
			},
		},
		{
			"invalid format fallback",
			"invalid-version-string",
			config.VersionString{
				Prefix: "",
				Major:  "invalid-version-string",
				Minor:  "",
				Patch:  "",
				Suffix: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.ParseVersion(tt.input)
			if result.Prefix != tt.expected.Prefix ||
				result.Major != tt.expected.Major ||
				result.Minor != tt.expected.Minor ||
				result.Patch != tt.expected.Patch ||
				result.Suffix != tt.expected.Suffix {
				t.Errorf("ParseVersion(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVersionString_Compare(t *testing.T) {
	tests := []struct {
		name     string
		a        config.VersionString
		b        config.VersionString
		expected int
	}{
		{
			"equal versions",
			config.VersionString{Prefix: "v", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "v", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			0,
		},
		{
			"different prefixes",
			config.VersionString{Prefix: "a", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "b", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			-1,
		},
		{
			"different major versions",
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "", Major: "2", Minor: "2", Patch: "3", Suffix: ""},
			-1,
		},
		{
			"different minor versions",
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "", Major: "1", Minor: "3", Patch: "3", Suffix: ""},
			-1,
		},
		{
			"different patch versions",
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "4", Suffix: ""},
			-1,
		},
		{
			"different suffixes",
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: "-alpha"},
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: "-beta"},
			-1,
		},
		{
			"no suffix vs suffix",
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: ""},
			config.VersionString{Prefix: "", Major: "1", Minor: "2", Patch: "3", Suffix: "-beta"},
			-1,
		},
		{
			"multi-digit comparison",
			config.VersionString{Prefix: "", Major: "10", Minor: "0", Patch: "0", Suffix: ""},
			config.VersionString{Prefix: "", Major: "2", Minor: "0", Patch: "0", Suffix: ""},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			if result != tt.expected {
				t.Errorf("VersionString.Compare() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestVersionString_Compare_Integration(t *testing.T) {
	versions := []string{
		"0.9.0",
		"1.0.0-alpha",
		"1.0.0-beta",
		"1.0.0",
		"1.0.1",
		"1.1.0",
		"2.0.0",
		"10.0.0",
	}

	for i := 0; i < len(versions)-1; i++ {
		for j := i + 1; j < len(versions); j++ {
			a := config.ParseVersion(versions[i])
			b := config.ParseVersion(versions[j])

			if a.Compare(b) >= 0 {
				t.Errorf("Expected %s < %s, but got %d", versions[i], versions[j], a.Compare(b))
			}

			if b.Compare(a) <= 0 {
				t.Errorf("Expected %s > %s, but got %d", versions[j], versions[i], b.Compare(a))
			}
		}
	}
}
