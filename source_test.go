package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSource_EnsureValid(t *testing.T) {
	t.Run("succeeds when directory exists", func(t *testing.T) {
		tmpDir := t.TempDir()

		source := NewSource(tmpDir)
		err := source.EnsureValid()
		require.NoError(t, err)
	})

	t.Run("fails when path does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, "does_not_exist")

		source := NewSource(nonExistentPath)
		err := source.EnsureValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "source path does not exist")
	})

	t.Run("fails when path is a file", func(t *testing.T) {
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "file.txt")

		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		source := NewSource(filePath)
		err = source.EnsureValid()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exists but is not a directory")
	})
}

func TestSource_RealPath(t *testing.T) {
	t.Run("joins paths correctly", func(t *testing.T) {
		source := NewSource("/base/path")
		result := source.RealPath("relative/file.txt")
		expected := filepath.Join("/base/path", "relative/file.txt")
		assert.Equal(t, expected, result)
	})

	t.Run("handles empty relative path", func(t *testing.T) {
		source := NewSource("/base/path")
		result := source.RealPath("")
		expected := filepath.Join("/base/path", "")
		assert.Equal(t, expected, result)
	})
}

func TestSource_Prepare(t *testing.T) {
	t.Run("processes files to copy", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files that should be copied
		testFiles := []string{
			"test.cpp",
			"config.json",
			"readme.txt",
			"model.p3d",
			"material.rvmat",
		}

		for _, file := range testFiles {
			filePath := filepath.Join(tmpDir, file)
			err := os.WriteFile(filePath, []byte("test content"), 0644)
			require.NoError(t, err)
		}

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)
		require.NotNil(t, task)

		// Verify all files are in the copy list
		assert.Len(t, task.Copy, len(testFiles))
		for _, file := range testFiles {
			assert.Contains(t, task.Copy, file)
		}

		// Verify manifest entries exist
		assert.Len(t, task.Manifest, len(testFiles))
		for _, file := range testFiles {
			hash, exists := task.Manifest[file]
			assert.True(t, exists, "manifest should contain %s", file)
			assert.Len(t, hash, 16, "hash should be 16 characters (8 bytes in hex)")
		}
	})

	t.Run("processes files to convert", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create test files that should be converted
		testFiles := []string{
			"image1.png",
			"image2.jpg",
			"photo.PNG", // Test case insensitivity
		}

		for _, file := range testFiles {
			filePath := filepath.Join(tmpDir, file)
			err := os.WriteFile(filePath, []byte("fake image data"), 0644)
			require.NoError(t, err)
		}

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)
		require.NotNil(t, task)

		// Verify all files are in the convert list
		assert.Len(t, task.Convert, len(testFiles))
		for _, file := range testFiles {
			assert.Contains(t, task.Convert, file)
		}

		// Verify manifest entries exist
		assert.Len(t, task.Manifest, len(testFiles))
	})

	t.Run("ignores unsupported file types", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create files that should be ignored
		ignoredFiles := []string{
			"script.sh",
			"data.xml",
			"archive.zip",
		}

		for _, file := range ignoredFiles {
			filePath := filepath.Join(tmpDir, file)
			err := os.WriteFile(filePath, []byte("content"), 0644)
			require.NoError(t, err)
		}

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)
		require.NotNil(t, task)

		// Verify ignored files are not in copy or convert lists
		assert.Len(t, task.Copy, 0)
		assert.Len(t, task.Convert, 0)
		assert.Len(t, task.Manifest, 0)
	})

	t.Run("handles nested directories", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create nested structure
		nestedPath := filepath.Join(tmpDir, "level1", "level2")
		err := os.MkdirAll(nestedPath, 0755)
		require.NoError(t, err)

		files := map[string]string{
			"root.cpp":                "root content",
			"level1/mid.json":         "mid content",
			"level1/level2/deep.txt":  "deep content",
			"level1/level2/image.png": "image data",
		}

		for file, content := range files {
			filePath := filepath.Join(tmpDir, file)
			err := os.MkdirAll(filepath.Dir(filePath), 0755)
			require.NoError(t, err)
			err = os.WriteFile(filePath, []byte(content), 0644)
			require.NoError(t, err)
		}

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)
		require.NotNil(t, task)

		// Should have 3 copy files and 1 convert file
		assert.Len(t, task.Copy, 3)
		assert.Len(t, task.Convert, 1)
		assert.Len(t, task.Manifest, 4)

		// Verify paths are relative and properly formatted
		assert.Contains(t, task.Copy, "root.cpp")
		assert.Contains(t, task.Copy, filepath.Join("level1", "mid.json"))
		assert.Contains(t, task.Copy, filepath.Join("level1", "level2", "deep.txt"))
		assert.Contains(t, task.Convert, filepath.Join("level1", "level2", "image.png"))
	})

	t.Run("empty directory produces empty task", func(t *testing.T) {
		tmpDir := t.TempDir()

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)
		require.NotNil(t, task)

		assert.Len(t, task.Copy, 0)
		assert.Len(t, task.Convert, 0)
		assert.Len(t, task.Manifest, 0)
	})

	t.Run("mixed file types processed correctly", func(t *testing.T) {
		tmpDir := t.TempDir()

		files := map[string]string{
			"code.cpp":   "copy",
			"data.json":  "copy",
			"image.png":  "convert",
			"photo.jpg":  "convert",
			"ignore.xml": "ignore",
			"skip.sh":    "ignore",
		}

		for file := range files {
			filePath := filepath.Join(tmpDir, file)
			err := os.WriteFile(filePath, []byte("content"), 0644)
			require.NoError(t, err)
		}

		source := NewSource(tmpDir)
		task, err := source.Prepare()
		require.NoError(t, err)

		assert.Len(t, task.Copy, 2, "should have 2 files to copy")
		assert.Len(t, task.Convert, 2, "should have 2 files to convert")
		assert.Len(t, task.Manifest, 4, "should have 4 files in manifest")

		// Verify the correct files are categorized
		assert.Contains(t, task.Copy, "code.cpp")
		assert.Contains(t, task.Copy, "data.json")
		assert.Contains(t, task.Convert, "image.png")
		assert.Contains(t, task.Convert, "photo.jpg")
	})
}
