package container

// Multistage generates a Dockerfile for a multistage build
func MultiStageBuild(currentContainer string, currentVersion string, toPath string, fromPath string, base string) string {

	// If toPath not defined, toPath is the same as fromPath
	if toPath == "" {
		toPath = fromPath
	}

	// Create Dockerfile base
	dockerfile := "FROM ghcr.io/autamus/" + currentContainer + ":" + currentVersion + " as base\n" +
		"FROM " + base + "\n" +
		"COPY --from=base " + fromPath + " " + toPath + "\n" +
		"ENV PATH=/opt/spack/bin:$PATH\n" +
		"WORKDIR /opt/spack\n" +
		"RUN rm -rf /opt/spack/.spack-db/\n" +
		"ENTRYPOINT [\"/bin/bash\"]\n"

	return dockerfile
}
