package gomod

import (
	"errors"
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func RelFilePathToModPath(p string) (string, error) {
	p, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	p = strings.ReplaceAll(p, "\\", "/")
	return GetDirectoryModuleName(p)
}

func GetDirectoryModuleName(dir string) (string, error) {
	if dir == "." {
		return "", errors.New("failed to find any go mod file")
	}
	base := path.Base(dir)
	parent := path.Dir(dir)
	found := false
	gomodFilePath := filepath.Join(dir, "go.mod")

	gomodPath := ""

	data, err := os.ReadFile(gomodFilePath)
	if err == nil {
		found = true
		gomodPath = modfile.ModulePath(data)
	}

	if !found {
		p, err := GetDirectoryModuleName(parent)
		if err != nil {
			return "", err
		}
		return p + "/" + base, nil
	}

	if gomodPath == "" {
		return "", errors.New("failed to get module path from a go.mod file")
	}

	return gomodPath, nil
}
