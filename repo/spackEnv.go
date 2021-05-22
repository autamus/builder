package repo

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SpackEnv struct {
	Spack Spack `yaml:"spack"`
}

type Spack struct {
	Specs     []string                 `yaml:"specs"`
	View      bool                     `yaml:"view"`
	Packages  map[string]SpackPackages `yaml:"packages"`
	Config    SpackConfig              `yaml:"config"`
	Container SpackContainer           `yaml:"container"`
}

type SpackPackages struct {
	Target []string `yaml:"target"`
}

type SpackConfig struct {
	Compiler                SpackConfigCompiler `yaml:"compiler"`
	InstallMissingCompilers bool                `yaml:"install_missing_compilers"`
}

type SpackConfigCompiler struct {
	Target string `yaml:"target"`
}

type SpackContainer struct {
	OSPackages SpackContainerPackages `yaml:"os_packages"`
	Strip      bool                   `yaml:"strip"`
}

type SpackContainerPackages struct {
	Build []string `yaml:"build"`
	Final []string `yaml:"final"`
}

func defaultEnv(defaultPath string) (output SpackEnv, err error) {
	input, err := ioutil.ReadFile(defaultPath)
	if err != nil {
		return output, err
	}
	err = yaml.Unmarshal(input, &output)
	return output, err
}

// ParseSpackEnv parses a spack environment into a go struct.
func ParseSpackEnv(defaultPath, containerPath string) (result SpackEnv, err error) {
	result, err = defaultEnv(defaultPath)
	if err != nil {
		return result, err
	}
	input, err := ioutil.ReadFile(containerPath)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(input, &result)
	return result, err
}
