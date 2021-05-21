package repo

import (
	"errors"
	"os"
	"path/filepath"
)

func FindPackagePath(packageName, packagesPath string) (packagePath string, err error) {
	// Search packages for specified package
	found := false
	err = filepath.Walk(packagesPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Base(filepath.Dir(path)) == packageName &&
			filepath.Base(path) == "package.py" {
			packagePath = path
			found = true
		}
		return nil
	})
	if !found {
		return packagePath, errors.New("package not found")
	}
	return packagePath, err
}
