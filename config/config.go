package config

import (
	"os"
	"reflect"
	"strings"
)

// Config defines the configuration struct for importing settings from ENV Variables
type Config struct {
	Containers containers
	General    general
	Multistage multistage
	Packages   packages
	Parsers    parsers
	Repository repository
}

type general struct {
	Version string
}

type packages struct {
	Path         string
	PublicKeyURL string
}

type containers struct {
	Path           string
	Current        string
	DefaultEnVPath string
}

type repository struct {
	Path          string
	DefaultBranch string
}

type multistage struct {
	Topath	string
	Frompath	string
	Base	string
	Suffix	string
}

type parsers struct {
	Loaded string
}

var (
	// Global is the configuration struct for the application.
	Global Config
)

func init() {
	defaultConfig()
	parseConfigEnv()
}

func defaultConfig() {
	Global.General.Version = "0.1.1"
	Global.Containers.Path = "containers/"
	Global.Containers.DefaultEnVPath = "default.yaml"
	Global.Packages.Path = "spack/"
	Global.Repository.Path = "."
	Global.Repository.DefaultBranch = "main"
	Global.Parsers.Loaded = "spack"
	Global.Multistage.Topath = ""
	Global.Multistage.Frompath = ""
	Global.Multistage.Base = "spack/ubuntu-bionic"
	Global.Multistage.Suffix = "-layers"
}

func parseConfigEnv() {
	numSubStructs := reflect.ValueOf(&Global).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(&Global).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			if fieldVal != "Version" {
				fieldName := structType.Field(j).Name
				for _, prefix := range []string{"BUILDER", "INPUT"} {
					evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
					evVal, evExists := os.LookupEnv(evName)
					if evExists && evVal != fieldVal {
						iter.FieldByName(fieldName).SetString(evVal)
					}
				}
			}
		}
	}
}
