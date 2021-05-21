package repo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const DockerType = "Dockerfile"
const SpackType = "Spack"

func GetContainerType(containersPath, containerName string) (result, resultPath string, err error) {
	dockerfile := false
	spackEnv := false
	// Search containers for specified container
	err = filepath.Walk(containersPath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "/"+containerName+"/") {
			match, _ := filepath.Match("spack.yaml", filepath.Base(path))
			if match {
				spackEnv = true
				resultPath = path
			}
			match, _ = filepath.Match("Dockerfile", filepath.Base(path))
			if match {
				dockerfile = true
				resultPath = path
			}
		}
		return nil
	})
	if err != nil {
		return result, resultPath, err
	}
	if dockerfile && spackEnv {
		return result, resultPath, errors.New("found both dockerfile and spack.yaml")
	}
	if dockerfile {
		return DockerType, resultPath, nil
	}
	return SpackType, resultPath, nil
}
