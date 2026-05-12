package main

import (
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
				return fmt.Errorf("error: could not create output directory %q: %w", o.path, err)
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

func (o *Output) Hash(path string) (string, error) {
	hash := fnv.New64a()

	f, err := os.Open(filepath.Join(o.path, path))
	if err != nil {
		return "", fmt.Errorf("error opening output file %q: %w", path, err)
	}
	defer f.Close()

	_, err = io.Copy(hash, f)
	if err != nil {
		return "", fmt.Errorf("error hashing output file %q: %w", path, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (o *Output) Convert(src, dst, imgToPaaPath string) (string, string, error) {
	dstFile := swapExtension(dst, ".paa")
	err := convertWithPath(
		src,
		filepath.Join(o.path, dstFile),
		imgToPaaPath,
	)
	if err != nil {
		return "", "", err
	}

	hash, err := o.Hash(dstFile)

	return dstFile, hash, err
}

func (o *Output) PathsToClean(task *Task) ([]string, error) {
	// files and directories to delete
	toClean := []string{}
	// a map of output files
	requiredOutputs := make(map[string]struct{}, len(task.Manifest))
	// a map of the directories required for those files
	requiredDirs := make(map[string]struct{}, len(task.Manifest))
	// directories which exist in the output root
	dirs := []string{}

	for _, entry := range task.Manifest {
		outputPath := entry.OutputPath
		if outputPath == "" {
			outputPath = entry.SourcePath
		}
		requiredOutputs[outputPath] = struct{}{}
		markAllDirectoriesInPathAsRequired(requiredDirs, filepath.Dir(outputPath))
	}

	err := fs.WalkDir(os.DirFS(o.path), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ignore manifest file and the root dir
		if path == ".build.manifest" || path == "." {
			return nil
		}

		if d.IsDir() {
			dirs = append(dirs, path)
			return nil
		}

		if _, ok := requiredOutputs[path]; ok {
			return nil
		}

		toClean = append(toClean, path)

		return nil
	})

	// sort the discovered directories by depth
	sort.Slice(dirs, func(i, j int) bool {
		iDepth := strings.Count(dirs[i], string(filepath.Separator))
		jDepth := strings.Count(dirs[j], string(filepath.Separator))
		if iDepth != jDepth {
			return iDepth > jDepth
		}
		return dirs[i] > dirs[j]
	})

	// check each directory against the list of retained paths and add
	// to cleanup if we do not need them for output files
	for _, path := range dirs {
		if _, ok := requiredDirs[path]; ok {
			continue
		}
		toClean = append(toClean, path)
	}

	return toClean, err
}

func (o *Output) Remove(path string) error {
	return os.Remove(filepath.Join(o.path, path))
}

func markAllDirectoriesInPathAsRequired(requiredDirs map[string]struct{}, path string) {
	for path != "." {
		requiredDirs[path] = struct{}{}
		path = filepath.Dir(path)
	}
}
