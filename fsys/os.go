package fsys

import (
	"os"
	"path/filepath"
)

func MakeAbsolute(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

func DirectoryMustExist(path string) {
	absPath := MakeAbsolute(path)
	info, err := os.Stat(absPath)
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic("not a directory: " + absPath)
	}
}
