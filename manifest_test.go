package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadManifest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Manifest
		wantErr bool
	}{
		{
			name:  "valid single entry",
			input: "path/to/file.txt\tabcdef0123456789\n",
			want: Manifest{
				"path/to/file.txt": "abcdef0123456789",
			},
			wantErr: false,
		},
		{
			name: "valid multiple entries",
			input: "file1.txt\t0123456789abcdef\n" +
				"path/to/file2.txt\tfedcba9876543210\n" +
				"another/file.dat\t1111222233334444\n",
			want: Manifest{
				"file1.txt":          "0123456789abcdef",
				"path/to/file2.txt":  "fedcba9876543210",
				"another/file.dat":   "1111222233334444",
			},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    Manifest{},
			wantErr: false,
		},
		{
			name:    "invalid format - missing tab",
			input:   "file.txt abcdef0123456789\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			input:   "file.txt\tabcdef0123456789\textra\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - too short",
			input:   "file.txt\tabcdef012345678\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - too long",
			input:   "file.txt\tabcdef01234567890\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - contains invalid characters",
			input:   "file.txt\tgbcdef0123456789\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - contains uppercase invalid character",
			input:   "file.txt\tGBCDEF0123456789\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "uppercase hex is valid",
			input: "file.txt\tABCDEF0123456789\n",
			want: Manifest{
				"file.txt": "ABCDEF0123456789",
			},
			wantErr: false,
		},
		{
			name:  "mixed case hex is valid",
			input: "file.txt\tAbCdEf0123456789\n",
			want: Manifest{
				"file.txt": "AbCdEf0123456789",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := LoadManifest(reader)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStoreManifest(t *testing.T) {
	tests := []struct {
		name     string
		manifest Manifest
		wantErr  bool
	}{
		{
			name: "single entry",
			manifest: Manifest{
				"path/to/file.txt": "abcdef0123456789",
			},
			wantErr: false,
		},
		{
			name: "multiple entries",
			manifest: Manifest{
				"file1.txt":         "0123456789abcdef",
				"path/to/file2.txt": "fedcba9876543210",
				"another/file.dat":  "1111222233334444",
			},
			wantErr: false,
		},
		{
			name:     "empty manifest",
			manifest: Manifest{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := StoreManifest(&buf, tt.manifest)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify we can read back what we wrote
			loaded, err := LoadManifest(&buf)
			require.NoError(t, err, "StoreManifest() produced invalid output")
			assert.Equal(t, tt.manifest, loaded)
		})
	}
}

func TestLoadStoreManifestRoundTrip(t *testing.T) {
	original := Manifest{
		"file1.txt":              "0123456789abcdef",
		"path/to/file2.txt":      "fedcba9876543210",
		"another/deeply/nested/": "aBcDeF0123456789",
		"file.dat":               "1111222233334444",
	}

	var buf bytes.Buffer
	err := StoreManifest(&buf, original)
	require.NoError(t, err)

	loaded, err := LoadManifest(&buf)
	require.NoError(t, err)

	assert.Equal(t, original, loaded)
}
