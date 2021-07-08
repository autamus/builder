package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	parser "github.com/autamus/binoc/repo"
	"github.com/autamus/builder/config"
	"github.com/autamus/builder/container"
	"github.com/autamus/builder/repo"
	"github.com/autamus/builder/spack"
)

func main() {
	fmt.Println()
	fmt.Print(` ____        _ _     _           
| __ ) _   _(_) | __| | ___ _ __ 
|  _ \| | | | | |/ _' |/ _ \ '__|
| |_) | |_| | | | (_| |  __/ |   
|____/ \__,_|_|_|\__,_|\___|_|   
`)
	fmt.Printf("Application Version: v%s\n", config.Global.General.Version)
	fmt.Println()
	// Initialize parser functionality
	parser.Init(strings.Split(config.Global.Parsers.Loaded, ","))

	// Set inital values for Repository
	path := config.Global.Repository.Path
	pubKeyURL := config.Global.Packages.PublicKeyURL
	packagesPath := filepath.Join(path, config.Global.Packages.Path)
	containersPath := filepath.Join(path, config.Global.Containers.Path)
	defaultEnvPath := filepath.Join(path, config.Global.Containers.DefaultEnVPath)

	// Declare container values
	currentContainer := config.Global.Containers.Current
	currentVersion := ""
	currentDockerfile := ""

	// Variables for multi-stage build container
	multistageToPath := config.Global.Multistage.Topath
	multistageFromPath := config.Global.Multistage.Frompath
	multistageBase := config.Global.Multistage.Base
	multistageSuffix := config.Global.Multistage.Suffix

	// Check if the current run is a PR
	prVal, prExists := os.LookupEnv("GITHUB_EVENT_NAME")
	isPR := prExists && prVal == "pull_request"

	// Get the type of the container from the repository.
	cType, cPath, err := repo.GetContainerType(containersPath, currentContainer)
	if err != nil {
		log.Fatal(err)
	}

	if cType == repo.SpackType {
		// If the container is a spack environment, find the main spec.
		spackEnv, err := repo.ParseSpackEnv(defaultEnvPath, cPath)
		if err != nil {
			log.Fatal(err)
		}
		// Find the path to the main spec package
		specPath, err := repo.FindPackagePath(spackEnv.Spack.Specs[0], packagesPath)
		if err != nil {
			log.Fatal(err)
		}
		// Parse package for main spec
		result, err := parser.Parse(specPath)
		if err != nil {
			log.Fatal(err)
		}
		// Set container version from package
		currentVersion = result.Package.GetLatestVersion().String()

		// Containerize SpackEnv to Dockerfile
		currentDockerfile, err = spack.Containerize(spackEnv, isPR, pubKeyURL)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		output, err := ioutil.ReadFile(cPath)
		if err != nil {
			log.Fatal(err)
		}
		currentDockerfile = string(output)
	}
	// Override an empty version with latest
	if currentVersion == "" {
		currentVersion = "latest"
	}

	// Write the Dockerfile out to Disk
	f, err := os.Create(filepath.Join(path, "Dockerfile"))
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(currentDockerfile)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	// Generate a multistage build with the container
	if multistageFromPath != "" {
		fmt.Printf("Generating multistage build container.\n")
		fmt.Println()
		multistageDockerfile := container.MultiStageBuild(currentContainer,
			currentVersion,
			multistageToPath,
			multistageFromPath,
			multistageBase)

		// Write the layer Dockerfile out to Disk
		f, err := os.Create(filepath.Join(path, "Dockerfile.multistage"))
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(multistageDockerfile)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()

		// Set multistage build suffix as output
		fmt.Printf("::set-output name=multistage_suffix::%s\n", multistageSuffix)
	}

	// Save Container Name and Version as Output
	fmt.Printf("::set-output name=container::%s\n", currentContainer)
	fmt.Printf("::set-output name=version::%s\n", currentVersion)
	fmt.Printf("::set-output name=type::%s\n", cType)
}
