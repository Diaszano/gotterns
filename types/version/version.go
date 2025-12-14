// Package version provides parsing, comparison, and string
// representation of semantic versions following the Semantic
// Versioning 2.0.0 specification.
//
// It includes utilities to safely parse version strings using TryParse,
// panic on invalid versions using Parse, and compare versions
// lexically and semantically.
package version

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// re defines the regular expression pattern that matches
// valid semantic version strings according to the SemVer 2.0.0 specification.
//
// Examples of valid versions:
//   - 1.0.0
//   - 1.0.0-alpha
//   - 1.0.0-alpha+001
//   - 1.0.0+build.123
var re = regexp.MustCompile(
	`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`,
)

// Common error variables for version validation.
var (
	// ErrBase is the base error used for all version-related errors.
	ErrBase = errors.New("version error")

	// ErrEmpty indicates that the provided version string is empty.
	ErrEmpty = errors.Join(ErrBase, errors.New("empty input"))

	// ErrInvalidFormat indicates that the provided version string
	// does not match the expected SemVer 2.0.0 format.
	ErrInvalidFormat = errors.Join(ErrBase, errors.New("invalid format"))
)

// Version represents a semantic version (SemVer 2.0.0).
// It includes the major, minor, and patch numbers,
// as well as optional prerelease and build metadata fields.
type Version struct {
	Major      uint64 `json:"major"`
	Minor      uint64 `json:"minor"`
	Patch      uint64 `json:"patch"`
	PreRelease string `json:"pre_release,omitempty"`
	Build      string `json:"build,omitempty"`
}

// String returns the semantic version as a formatted string
// without a leading 'v' prefix.
//
// Example outputs:
//   - "1.0.0"
//   - "1.0.0-alpha"
//   - "1.0.0-alpha+build.001"
func (v Version) String() string {
	version := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.PreRelease != "" {
		version += fmt.Sprintf("-%s", v.PreRelease)
	}
	if v.Build != "" {
		version += fmt.Sprintf("+%s", v.Build)
	}

	return version
}

// Version returns the version string prefixed with 'v',
// commonly used in Git tags or Go module versions.
//
// Example outputs:
//   - "v1.0.0"
//   - "v1.0.0-beta"
//   - "v1.0.0+build.45"
func (v Version) Version() string {
	return "v" + v.String()
}

// Compare compares the current version with another Version.
//
// Returns:
//   - 1 if v > other
//   - -1 if v < other
//   - 0 if v == other
//
// Comparison is done in the following order:
//  1. Major
//  2. Minor
//  3. Patch
//  4. PreRelease (pre-release versions are considered lower than normal versions)
//
// Note: Build metadata is ignored in comparisons.
func (v Version) Compare(other Version) int8 {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	// Handle pre-release comparison
	if v.PreRelease == "" && other.PreRelease != "" {
		// v is a release, other is pre-release → v > other
		return 1
	}

	if v.PreRelease != "" && other.PreRelease == "" {
		// v is pre-release, other is release → v < other
		return -1
	}

	if v.PreRelease != "" && other.PreRelease != "" {
		if v.PreRelease > other.PreRelease {
			return 1
		} else if v.PreRelease < other.PreRelease {
			return -1
		}
	}

	// Versions are equal (including pre-release)
	return 0
}

// TryParse attempts to parse a semantic version string into a Version struct.
//
// It safely returns an error instead of panicking when the input is invalid.
// The function supports both versions with and without a 'v' prefix.
//
// Examples:
//
//	version, err := TryParse("v1.2.3")
//	version, err := TryParse("1.2.3-beta+build.12")
//
// Possible errors:
//   - ErrEmpty if the input string is empty
//   - ErrInvalidFormat if the version does not conform to the SemVer specification
func TryParse(input string) (Version, error) {
	if input == "" {
		return Version{}, ErrEmpty
	}

	version := strings.TrimPrefix(input, "v")

	if !re.MatchString(version) {
		return Version{}, ErrInvalidFormat
	}

	matches := re.FindStringSubmatch(version)

	if len(matches) < 6 {
		return Version{}, ErrInvalidFormat
	}

	major, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return Version{}, errors.Join(ErrBase, err)
	}

	minor, err := strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		return Version{}, errors.Join(ErrBase, err)
	}

	patch, err := strconv.ParseUint(matches[3], 10, 64)
	if err != nil {
		return Version{}, errors.Join(ErrBase, err)
	}

	return Version{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		PreRelease: matches[4],
		Build:      matches[5],
	}, nil
}

// Parse parses a semantic version string and returns a Version struct.
//
// Unlike TryParse, Parse panics if the input is empty or invalid.
// This function is suitable when invalid input indicates a programmer error.
//
// Valid examples:
//   - "1.0.0"
//   - "v1.2.3-alpha+build.45"
//
// Panics with:
//   - ErrEmpty if the version string is empty.
//   - ErrInvalidFormat if the version string has an invalid format.
func Parse(input string) Version {
	version, err := TryParse(input)
	if err != nil {
		panic(err)
	}

	return version
}
