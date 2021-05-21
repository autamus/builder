package repo

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type SpackEnv struct {
	Spack Spack `yaml:"spack"`
}

type Spack struct {
	Specs     []string                 `yaml:"specs,omitempty"`
	View      bool                     `yaml:"view,omitempty"`
	Packages  map[string]SpackPackages `yaml:"packages,omitempty"`
	Config    SpackConfig              `yaml:"config,omitempty"`
	Container SpackContainer           `yaml:"container,omitempty"`
}

type SpackPackages struct {
	Target []string `yaml:"target,omitempty"`
}

type SpackConfig struct {
	Compiler                SpackConfigCompiler `yaml:"compiler,omitempty"`
	InstallMissingCompilers bool                `yaml:"install_missing_compilers,omitempty"`
}

type SpackConfigCompiler struct {
	Target string `yaml:"target,omitempty"`
}

type SpackContainer struct {
	OSPackages SpackContainerPackages `yaml:"os_packages,omitempty"`
	Strip      bool                   `yaml:"strip,omitempty"`
}

type SpackContainerPackages struct {
	Build []string `yaml:"build,omitempty"`
	Final []string `yaml:"final,omitempty"`
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
