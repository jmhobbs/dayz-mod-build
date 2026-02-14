package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Manifest map[string]string

var hashRegexp = regexp.MustCompile(`^[a-fA-F0-9]{16}$`)

func LoadManifest(in io.Reader) (Manifest, error) {
	manifest := make(Manifest)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid manifest line: %s", scanner.Text())
		}

		if !hashRegexp.MatchString(parts[1]) {
			return nil, fmt.Errorf("invalid hash in manifest line: %s", scanner.Text())
		}

		manifest[parts[0]] = parts[1]
	}
	return manifest, scanner.Err()
}

func StoreManifest(out io.Writer, manifest Manifest) error {
	for path, hash := range manifest {
		_, err := fmt.Fprintf(out, "%s\t%s\n", path, hash)
		if err != nil {
			return err
		}
	}
	return nil
}
