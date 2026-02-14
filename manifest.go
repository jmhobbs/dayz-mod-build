package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type ManifestEntry struct {
	SourcePath string
	SourceHash string
	OutputPath string
	OutputHash string
}

type Manifest map[string]ManifestEntry

var hashRegexp = regexp.MustCompile(`^[a-fA-F0-9]{16}$`)

func LoadManifest(in io.Reader) (Manifest, error) {
	manifest := make(Manifest)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) == 2 {
			if !hashRegexp.MatchString(parts[1]) {
				return nil, fmt.Errorf("invalid hash in manifest line: %s", scanner.Text())
			}
			manifest[parts[0]] = ManifestEntry{
				SourcePath: parts[0],
				SourceHash: parts[1],
				OutputPath: parts[0],
				OutputHash: parts[1],
			}

		} else if len(parts) == 4 {
			if !hashRegexp.MatchString(parts[1]) || !hashRegexp.MatchString(parts[3]) {
				return nil, fmt.Errorf("invalid hash in manifest line: %s", scanner.Text())
			}
			manifest[parts[0]] = ManifestEntry{
				SourcePath: parts[0],
				SourceHash: parts[1],
				OutputPath: parts[2],
				OutputHash: parts[3],
			}
		} else {
			return nil, fmt.Errorf("invalid manifest line: %s", scanner.Text())
		}

	}
	return manifest, scanner.Err()
}

func StoreManifest(out io.Writer, manifest Manifest) error {
	var err error
	for _, entry := range manifest {
		if entry.OutputPath != "" && entry.OutputHash != "" {
			_, err = fmt.Fprintf(out, "%s\t%s\t%s\t%s\n", entry.SourcePath, entry.SourceHash, entry.OutputPath, entry.OutputHash)
		} else {
			_, err = fmt.Fprintf(out, "%s\t%s\n", entry.SourcePath, entry.SourceHash)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
