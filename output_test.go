package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutput_EnsureExists(t *testing.T) {
	t.Run("creates directory when it does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "new_output_dir")

		output := NewOutput(testPath)
		err := output.EnsureExists()
		require.NoError(t, err)

		// Verify directory was created
		finfo, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, finfo.IsDir())
	})

	t.Run("creates nested directories when they do not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "level1", "level2", "level3")

		output := NewOutput(testPath)
		err := output.EnsureExists()
		require.NoError(t, err)

		// Verify nested directories were created
		finfo, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, finfo.IsDir())
	})

	t.Run("succeeds when directory already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "existing_dir")

		// Create directory first
		err := os.Mkdir(testPath, 0755)
		require.NoError(t, err)

		output := NewOutput(testPath)
		err = output.EnsureExists()
		require.NoError(t, err)

		// Verify it's still a directory
		finfo, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.True(t, finfo.IsDir())
	})

	t.Run("fails when path exists but is a file", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "file_not_dir")

		// Create a file at the path
		err := os.WriteFile(testPath, []byte("test"), 0644)
		require.NoError(t, err)

		output := NewOutput(testPath)
		err = output.EnsureExists()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "exists but is not a directory")
	})

	t.Run("creates directory with correct permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		testPath := filepath.Join(tmpDir, "perm_test")

		output := NewOutput(testPath)
		err := output.EnsureExists()
		require.NoError(t, err)

		// Verify permissions (0755)
		finfo, err := os.Stat(testPath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0755)|os.ModeDir, finfo.Mode())
	})

	t.Run("fails with invalid parent path", func(t *testing.T) {
		tmpDir := t.TempDir()
		fileInPath := filepath.Join(tmpDir, "file_in_path")
		testPath := filepath.Join(fileInPath, "cannot_create")

		// Create a file where a parent directory should be
		err := os.WriteFile(fileInPath, []byte("blocking"), 0644)
		require.NoError(t, err)

		output := NewOutput(testPath)
		err = output.EnsureExists()
		require.Error(t, err)
	})
}

func TestOutput_PathsToClean(t *testing.T) {
	t.Run("does not schedule non-empty ancestors of retained copied files", func(t *testing.T) {
		tmpDir := t.TempDir()
		output := NewOutput(tmpDir)

		keepPath := filepath.Join("level1", "level2", "keep.txt")
		deletePath := filepath.Join("level1", "delete.txt")

		writeOutputTestFile(t, tmpDir, keepPath, "keep")
		writeOutputTestFile(t, tmpDir, deletePath, "delete")

		task := &Task{
			Manifest: Manifest{
				keepPath: {
					SourcePath: keepPath,
					SourceHash: "0123456789abcdef",
				},
			},
		}

		toClean, err := output.PathsToClean(task)
		require.NoError(t, err)

		assert.Contains(t, toClean, deletePath)
		assert.NotContains(t, toClean, filepath.Join("level1"))
		assert.NotContains(t, toClean, ".")
	})

	t.Run("does not schedule non-empty ancestors of retained converted files", func(t *testing.T) {
		tmpDir := t.TempDir()
		output := NewOutput(tmpDir)

		keepSourcePath := filepath.Join("textures", "nested", "keep.png")
		keepOutputPath := swapExtension(keepSourcePath, ".paa")
		deletePath := filepath.Join("textures", "old.paa")

		writeOutputTestFile(t, tmpDir, keepOutputPath, "keep")
		writeOutputTestFile(t, tmpDir, deletePath, "delete")

		task := &Task{
			Manifest: Manifest{
				keepSourcePath: {
					SourcePath: keepSourcePath,
					SourceHash: "fedcba9876543210",
					OutputPath: keepOutputPath,
					OutputHash: "0011223344556677",
				},
			},
		}

		toClean, err := output.PathsToClean(task)
		require.NoError(t, err)

		assert.Contains(t, toClean, deletePath)
		assert.NotContains(t, toClean, filepath.Join("textures"))
		assert.NotContains(t, toClean, ".")
	})
}

func writeOutputTestFile(t *testing.T, root, relativePath, contents string) {
	t.Helper()

	fullPath := filepath.Join(root, relativePath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	require.NoError(t, err)

	err = os.WriteFile(fullPath, []byte(contents), 0644)
	require.NoError(t, err)
}
