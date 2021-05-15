package gomod

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"io/fs"
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
	gomodPath := ""

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		if info.Name() == "go.mod" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			found = true

			gomodPath = modfile.ModulePath(data)

			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return "", err
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
