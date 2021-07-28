package repo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// FindPackagePath returns the path to the package.py file for a given package.
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

// GetPackageFromSpec extracts the name of the package from the abstract spec.
func GetPackageFromSpec(spec string) (packageName string) {
	sep := strings.IndexFunc(spec, func(token rune) bool {
		specials := []rune{'%', '^', '~', '+', '@'}
		for _, special := range specials {
			if token == special {
				return true
			}
		}
		return false
	})
	if sep > 0 {
		packageName = spec[:sep]
	} else {
		packageName = spec
	}
	return packageName
}
