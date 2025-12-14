package version_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Diaszano/gotterns/types/version"
)

func TestTryParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    version.Version
		wantErr error
	}{
		{
			name:  "valid semantic version",
			input: "1.0.0",
			want: version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			wantErr: nil,
		},
		{
			name:  "version with v prefix",
			input: "v2.1.3",
			want: version.Version{
				Major: 2,
				Minor: 1,
				Patch: 3,
			},
			wantErr: nil,
		},
		{
			name:  "version with prerelease",
			input: "1.2.3-alpha",
			want: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
			},
			wantErr: nil,
		},
		{
			name:  "version with build metadata",
			input: "1.2.3+build.123",
			want: version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
			wantErr: nil,
		},
		{
			name:  "version with prerelease and build metadata",
			input: "1.2.3-beta+exp.sha.5114f85",
			want: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "beta",
				Build:      "exp.sha.5114f85",
			},
			wantErr: nil,
		},
		{
			name:    "empty input returns ErrEmpty",
			input:   "",
			wantErr: version.ErrEmpty,
		},
		{
			name:    "invalid format returns ErrInvalidFormat",
			input:   "invalid-version",
			wantErr: version.ErrInvalidFormat,
		},
		{
			name:    "missing patch number returns ErrInvalidFormat",
			input:   "1.2",
			wantErr: version.ErrInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := version.TryParse(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err, "expected error for input %q", tt.input)
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tt.want, got, "expected nil version for invalid input")
				return
			}

			require.NoError(t, err, "unexpected error for input %q", tt.input)
			require.Equal(t, tt.want, got, "unexpected parsed version for input %q", tt.input)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        version.Version
		expectPanic bool
		err         error
	}{
		{
			name:  "valid semantic version",
			input: "1.2.3",
			want: version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:  "version with prefix v",
			input: "v1.0.0",
			want: version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
		{
			name:  "version with prerelease",
			input: "1.2.3-alpha",
			want: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
			},
		},
		{
			name:  "version with build metadata",
			input: "1.2.3+build.123",
			want: version.Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
		},
		{
			name:  "version with prerelease and build metadata",
			input: "1.2.3-beta+exp.sha.5114f85",
			want: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "beta",
				Build:      "exp.sha.5114f85",
			},
		},
		{
			name:        "empty input should panic",
			input:       "",
			expectPanic: true,
			err:         version.ErrEmpty,
		},
		{
			name:        "invalid format should panic",
			input:       "invalid-version",
			expectPanic: true,
			err:         version.ErrInvalidFormat,
		},
		{
			name:        "missing patch number should panic",
			input:       "1.2",
			expectPanic: true,
			err:         version.ErrInvalidFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				require.PanicsWithError(t, tt.err.Error(), func() {
					version.Parse(tt.input)
				}, "expected panic for input %q", tt.input)
				return
			}

			got := version.Parse(tt.input)
			require.Equal(t, tt.want, got, "unexpected version struct for input %q", tt.input)
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name string
		v    version.Version
		want string
	}{
		{
			name: "basic version",
			v:    version.Version{Major: 1, Minor: 2, Patch: 3},
			want: "1.2.3",
		},
		{
			name: "version with prerelease",
			v:    version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "beta"},
			want: "1.2.3-beta",
		},
		{
			name: "version with build metadata",
			v:    version.Version{Major: 1, Minor: 2, Patch: 3, Build: "build.123"},
			want: "1.2.3+build.123",
		},
		{
			name: "version with prerelease and build metadata",
			v: version.Version{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha",
				Build:      "exp.sha",
			},
			want: "1.2.3-alpha+exp.sha",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.String()
			require.Equal(t, tt.want, got, "unexpected string output")
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name     string
		v1       version.Version
		v2       version.Version
		expected int8
	}{
		{
			name:     "major greater",
			v1:       version.Version{Major: 2, Minor: 0, Patch: 0},
			v2:       version.Version{Major: 1, Minor: 5, Patch: 9},
			expected: 1,
		},
		{
			name:     "major smaller",
			v1:       version.Version{Major: 1, Minor: 0, Patch: 0},
			v2:       version.Version{Major: 2, Minor: 0, Patch: 0},
			expected: -1,
		},
		{
			name:     "minor greater",
			v1:       version.Version{Major: 1, Minor: 3, Patch: 0},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 9},
			expected: 1,
		},
		{
			name:     "minor smaller",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 0},
			v2:       version.Version{Major: 1, Minor: 3, Patch: 0},
			expected: -1,
		},
		{
			name:     "patch greater",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 5},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3},
			expected: 1,
		},
		{
			name:     "patch smaller",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 5},
			expected: -1,
		},
		{
			name:     "pre-release vs release",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3},
			expected: -1,
		},
		{
			name:     "release vs pre-release",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "beta"},
			expected: 1,
		},
		{
			name:     "pre-release greater",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "beta"},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"},
			expected: 1,
		},
		{
			name:     "pre-release smaller",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "beta"},
			expected: -1,
		},
		{
			name:     "versions equal",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3},
			expected: 0,
		},
		{
			name:     "pre-release equal",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3, PreRelease: "alpha"},
			expected: 0,
		},
		{
			name:     "build ignored",
			v1:       version.Version{Major: 1, Minor: 2, Patch: 3, Build: "001"},
			v2:       version.Version{Major: 1, Minor: 2, Patch: 3, Build: "002"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.v1.Compare(tt.v2)
			require.Equal(t, tt.expected, result)
		})
	}
}
