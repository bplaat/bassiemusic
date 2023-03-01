package utils

import (
	"os"
	"path/filepath"
)

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !file.IsDir() {
			size += file.Size()
		}
		return err
	})
	return size, err
}
