package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func BuildRootExists(outputRoot string) (bool, error) {
	finfo, err := os.Stat(outputRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !finfo.IsDir() {
		return false, fmt.Errorf("error: build output path %q exists but is not a directory", outputRoot)
	}
	return true, nil
}

func GetAddonsToBuild(sourceRoot string) ([]string, error) {
	addonsToBuild := []string{}

	entries, err := os.ReadDir(sourceRoot)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			addonsToBuild = append(addonsToBuild, entry.Name())
		}
	}

	return addonsToBuild, nil
}

func GetAddonOutputDirsToClean(outputRoot string, addons []string) ([]string, error) {
	addonsToClean := []string{}
	for _, addon := range addons {
		_, err := os.Stat(filepath.Join(outputRoot, addon))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		addonsToClean = append(addonsToClean, addon)
	}
	return addonsToClean, nil
}
