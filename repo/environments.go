package repo

import "strings"

func GetArchString(in SpackEnv) string {
	arches := make(map[string]bool)

	for _, target := range in.Spack.Config.Compiler.Target {
		switch target {
		case "x86_64_v3":
			arches["linux/amd64"] = true
		case "aarch64":
			arches["linux/arm64"] = true
		}
	}

	result := []string{}

	for key := range arches {
		result = append(result, key)
	}

	return strings.Join(result, ",")
}
