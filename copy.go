package main

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var formatsToCopy = []string{
	".c",
	".cpp",
	".csv",
	".edds",
	".emat",
	".fnt",
	".html",
	".imageset",
	".json",
	".layout",
	".map",
	".ogg",
	".p3d",
	".paa",
	".ptc",
	".rvmat",
	".tga",
	".txt",
	".xml",
}

func shouldCopy(path string) bool {
	return slices.Contains(formatsToCopy, strings.ToLower(filepath.Ext(path)))
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	sink, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = io.Copy(sink, source)
	return err
}

func copyFileWithPath(src, dst string) error {
	err := os.MkdirAll(filepath.Dir(dst), 0755)
	if err != nil {
		return err
	}
	return copyFile(src, dst)
}
