package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var formatsToConvert = []string{".png", ".jpg"}

func swapExtension(path, newExt string) string {
	stem := strings.TrimSuffix(path, filepath.Ext(path))
	return stem + newExt
}

func shouldConvert(path string) bool {
	return slices.Contains(formatsToConvert, strings.ToLower(filepath.Ext(path)))
}

func convertWithPath(src, dst, converter string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}
	cmd := exec.Command(converter, src, dst)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error converting image: %w\n%s", err, string(out))
	}

	return nil
}
