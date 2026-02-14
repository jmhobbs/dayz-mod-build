package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Output struct {
	path string
}

func NewOutput(path string) *Output {
	return &Output{path: path}
}

func (o *Output) EnsureExists() error {
	finfo, err := os.Stat(o.path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(o.path, 0755)
			if err != nil {
				return fmt.Errorf("error: could not create output directory %q: %v", o.path, err)
			}
			return nil
		}
		return err
	}
	if !finfo.IsDir() {
		return fmt.Errorf("error: path %q exists but is not a directory", o.path)
	}
	return nil
}

func (o *Output) LoadManifest() (Manifest, error) {
	f, err := os.Open(filepath.Join(o.path, ".build.manifest"))
	if err != nil {
		if os.IsNotExist(err) {
			return make(Manifest), nil
		}
	}
	defer f.Close()
	return LoadManifest(f)
}

func (o *Output) WriteManifest(manifest Manifest) error {
	f, err := os.Create(filepath.Join(o.path, ".build.manifest"))
	if err != nil {
		return err
	}
	defer f.Close()
	return StoreManifest(f, manifest)
}

func (o *Output) Copy(src, dst string) error {
	return copyFileWithPath(src, filepath.Join(o.path, dst))
}

func (o *Output) Convert(src, dst, imgToPaaPath string) error {
	return convertWithPath(
		src,
		filepath.Join(o.path, swapExtension(dst, ".paa")),
		imgToPaaPath,
	)
}

func (o *Output) PathsToClean(task *Task) ([]string, error) {
	toClean := []string{}
	return toClean, fs.WalkDir(os.DirFS(o.path), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == path && path == ".build.manifest" {
			return nil
		}

		if d.IsDir() {
			// TODO: Remove empty directories
			return nil
		}

		// straight copy
		if _, ok := task.Manifest[path]; ok {
			return nil
		}

		// conversions
		possibleSources := []string{
			swapExtension(path, ".png"),
			swapExtension(path, ".jpg"),
		}

		for _, possiblePath := range possibleSources {
			if _, ok := task.Manifest[possiblePath]; ok {
				return nil
			}
		}

		toClean = append(toClean, path)

		return nil
	})
}

func (o *Output) Remove(path string) error {
	return os.Remove(filepath.Join(o.path, path))
}
