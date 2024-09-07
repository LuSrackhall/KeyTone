package audioPackageList

import (
	"os"
	"path/filepath"
)

func GetAudioPackageList(rootDir string) ([]string, error) {
	var directories []string

	entries, err := os.ReadDir(rootDir)

	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(rootDir, entry.Name()))
		}
	}

	return directories, err
}
