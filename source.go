package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Source struct {
	path string
}

func NewSource(path string) *Source {
	return &Source{path: path}
}

func (s *Source) EnsureValid() error {
	finfo, err := os.Stat(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("source path does not exist: %q", s.path)
		}
		return err
	}
	if !finfo.IsDir() {
		return fmt.Errorf("error: source path %q exists but is not a directory", s.path)
	}
	return nil
}

type Task struct {
	Manifest Manifest
	Copy     []string
	Convert  []string
}

func (s *Source) Prepare() (*Task, error) {
	task := &Task{Manifest: make(Manifest), Copy: []string{}}

	hash := fnv.New64a()

	err := fs.WalkDir(os.DirFS(s.path), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		shouldManifest := false
		if shouldCopy(path) {
			task.Copy = append(task.Copy, path)
			shouldManifest = true
		} else if shouldConvert(path) {
			task.Convert = append(task.Convert, path)
			shouldManifest = true
		}

		if !shouldManifest {
			return nil
		}

		f, err := os.Open(s.RealPath(path))
		if err != nil {
			return fmt.Errorf("error opening input file %q: %v", path, err)
		}
		defer f.Close()

		hash.Reset()
		_, err = io.Copy(hash, f)
		if err != nil {
			return fmt.Errorf("error hashing input file %q: %v", path, err)
		}
		task.Manifest[path] = fmt.Sprintf("%x", hash.Sum(nil))

		return nil
	})

	return task, err
}

func (s *Source) RealPath(path string) string {
	return filepath.Join(s.path, path)
}
