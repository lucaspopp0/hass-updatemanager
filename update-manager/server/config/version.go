package config

import (
	"regexp"
	"strconv"
)

type VersionSegment string

func (a VersionSegment) Compare(b VersionSegment) int {
	if a == b {
		return 0
	}

	intA, errA := strconv.Atoi(string(a))
	intB, errB := strconv.Atoi(string(b))

	if errA == nil && errB == nil {
		if intA < intB {
			return -1
		}

		return 1
	}

	if a < b {
		return -1
	}

	return 1
}

type VersionString struct {
	Prefix string

	Major VersionSegment
	Minor VersionSegment
	Patch VersionSegment

	Suffix string
}

func ParseVersion(rawVersion string) VersionString {
	version := VersionString{}

	expr := regexp.MustCompile(
		`^(?<prefix>\D*)(?<major>\d+)\.?(?<minor>\w*)\.?(?<patch>\w*)(?<suffix>(?:\W.*|))$`,
	)

	match := expr.FindStringSubmatch(rawVersion)

	if len(match) == 6 {
		version.Prefix = match[1]
		version.Major = VersionSegment(match[2])
		version.Minor = VersionSegment(match[3])
		version.Patch = VersionSegment(match[4])
		version.Suffix = match[5]
	} else {
		version.Major = VersionSegment(rawVersion)
	}

	return version
}

func (a VersionString) Compare(b VersionString) int {
	switch {
	case a.Prefix < b.Prefix:
		return -1
	case a.Prefix > b.Prefix:
		return 1
	}

	if majors := a.Major.Compare(b.Major); majors != 0 {
		return majors
	} else if minors := a.Minor.Compare(b.Minor); minors != 0 {
		return minors
	} else if patches := a.Patch.Compare(b.Patch); patches != 0 {
		return patches
	} else {
		switch {
		case a.Suffix == b.Suffix:
			return 0
		case a.Suffix < b.Suffix:
			return -1
		default:
			return 1
		}
	}
}
