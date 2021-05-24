package repo

import (
	"io/ioutil"
	"strings"

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
		return
	}
	input, err := ioutil.ReadFile(containerPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(input, &result)
	if err != nil {
		return
	}

	// Clean specs from varient/version information.
	specs := make([]string, len(result.Spack.Specs))
	for _, spec := range result.Spack.Specs {
		i := strings.IndexFunc(spec, versend)
		if i > 0 {
			specs = append(specs, spec[:i])
		} else {
			specs = append(specs, spec)
		}
	}
	result.Spack.Specs = specs

	return result, nil
}

// versend returns true at the end of the name of a dependency
func versend(input rune) bool {
	for _, c := range []rune{'@', '~', '+'} {
		if input == c {
			return true
		}
	}
	return false
}
