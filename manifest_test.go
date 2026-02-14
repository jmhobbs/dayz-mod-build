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
			name:  "valid single entry (2 parts)",
			input: "path/to/file.txt\tabcdef0123456789\n",
			want: Manifest{
				"path/to/file.txt": ManifestEntry{
					SourcePath: "path/to/file.txt",
					SourceHash: "abcdef0123456789",
					OutputPath: "path/to/file.txt",
					OutputHash: "abcdef0123456789",
				},
			},
			wantErr: false,
		},
		{
			name:  "valid single entry (4 parts)",
			input: "file_co.png\tabcdef0123456789\tfile_co.paa\t1234567890abcdef\n",
			want: Manifest{
				"file_co.png": ManifestEntry{
					SourcePath: "file_co.png",
					SourceHash: "abcdef0123456789",
					OutputPath: "file_co.paa",
					OutputHash: "1234567890abcdef",
				},
			},
			wantErr: false,
		},
		{
			name: "valid multiple entries",
			input: "file1.txt\t0123456789abcdef\n" +
				"path/to/file2.txt\tfedcba9876543210\n" +
				"another/file.cpp\t1111222233334444\n",
			want: Manifest{
				"file1.txt": ManifestEntry{
					SourcePath: "file1.txt",
					SourceHash: "0123456789abcdef",
					OutputPath: "file1.txt",
					OutputHash: "0123456789abcdef",
				},
				"path/to/file2.txt": ManifestEntry{
					SourcePath: "path/to/file2.txt",
					SourceHash: "fedcba9876543210",
					OutputPath: "path/to/file2.txt",
					OutputHash: "fedcba9876543210",
				},
				"another/file.cpp": ManifestEntry{
					SourcePath: "another/file.cpp",
					SourceHash: "1111222233334444",
					OutputPath: "another/file.cpp",
					OutputHash: "1111222233334444",
				},
			},
			wantErr: false,
		},
		{
			name: "valid mixed 2 and 4 part entries",
			input: "file1.txt\t0123456789abcdef\n" +
				"image.png\tabcdef0123456789\timage.paa\t9876543210fedcba\n",
			want: Manifest{
				"file1.txt": ManifestEntry{
					SourcePath: "file1.txt",
					SourceHash: "0123456789abcdef",
					OutputPath: "file1.txt",
					OutputHash: "0123456789abcdef",
				},
				"image.png": ManifestEntry{
					SourcePath: "image.png",
					SourceHash: "abcdef0123456789",
					OutputPath: "image.paa",
					OutputHash: "9876543210fedcba",
				},
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
			name:    "invalid format - 3 parts",
			input:   "file.txt\tabcdef0123456789\textra\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid format - 5 parts",
			input:   "file.txt\tabcdef0123456789\textra\t1234567890abcdef\tmore\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - too short (2 part)",
			input:   "file.txt\tabcdef012345678\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - too long (2 part)",
			input:   "file.txt\tabcdef01234567890\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - source hash invalid (4 part)",
			input:   "file.txt\tgbcdef0123456789\toutput.txt\t1234567890abcdef\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hash - output hash invalid (4 part)",
			input:   "file.txt\tabcdef0123456789\toutput.txt\tgbcdef0123456789\n",
			want:    nil,
			wantErr: true,
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
			name: "single entry (2 part format)",
			manifest: Manifest{
				"path/to/file.txt": ManifestEntry{
					SourcePath: "path/to/file.txt",
					SourceHash: "abcdef0123456789",
					OutputPath: "path/to/file.txt",
					OutputHash: "abcdef0123456789",
				},
			},
			wantErr: false,
		},
		{
			name: "single entry (4 part format)",
			manifest: Manifest{
				"file_co.png": ManifestEntry{
					SourcePath: "file_co.png",
					SourceHash: "abcdef0123456789",
					OutputPath: "file_co.paa",
					OutputHash: "1234567890abcdef",
				},
			},
			wantErr: false,
		},
		{
			name: "multiple entries mixed format",
			manifest: Manifest{
				"file1.txt": ManifestEntry{
					SourcePath: "file1.txt",
					SourceHash: "0123456789abcdef",
					OutputPath: "file1.txt",
					OutputHash: "0123456789abcdef",
				},
				"path/to/file2.txt": ManifestEntry{
					SourcePath: "path/to/file2.txt",
					SourceHash: "fedcba9876543210",
					OutputPath: "path/to/file2.txt",
					OutputHash: "fedcba9876543210",
				},
				"image.png": ManifestEntry{
					SourcePath: "image.png",
					SourceHash: "1111222233334444",
					OutputPath: "image.paa",
					OutputHash: "5555666677778888",
				},
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
		"file1.txt": ManifestEntry{
			SourcePath: "file1.txt",
			SourceHash: "0123456789abcdef",
			OutputPath: "file1.txt",
			OutputHash: "0123456789abcdef",
		},
		"path/to/file2.txt": ManifestEntry{
			SourcePath: "path/to/file2.txt",
			SourceHash: "fedcba9876543210",
			OutputPath: "path/to/file2.txt",
			OutputHash: "fedcba9876543210",
		},
		"image.png": ManifestEntry{
			SourcePath: "image.png",
			SourceHash: "abcdef0123456789",
			OutputPath: "image.paa",
			OutputHash: "1111222233334444",
		},
	}

	var buf bytes.Buffer
	err := StoreManifest(&buf, original)
	require.NoError(t, err)

	loaded, err := LoadManifest(&buf)
	require.NoError(t, err)

	assert.Equal(t, original, loaded)
}
