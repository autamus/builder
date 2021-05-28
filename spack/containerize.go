package spack

import (
	"os"
	"os/exec"
	"strings"

	"github.com/autamus/builder/repo"
	"gopkg.in/yaml.v2"
)

// Containerize runs spack containerize on an environment
// and returns the resulting dockerfile.
func Containerize(sEnv repo.SpackEnv, isPR bool) (dockerfile string, err error) {
	// Create spack.yaml file.
	f, err := os.Create("spack.yaml")
	if err != nil {
		return
	}
	// Write data out to file.
	envFile, err := yaml.Marshal(sEnv)
	if err != nil {
		return
	}
	_, err = f.Write(envFile)
	if err != nil {
		return
	}
	f.Close()

	// Run Spack Containerize on File
	cmd := exec.Command("spack", "containerize")
	out, err := cmd.Output()
	if err != nil {
		return
	}
	dockerfile = string(out)

	// Delete File
	err = os.Remove("spack.yaml")
	if err != nil {
		return
	}

	// Add Docker ENV to dockerfile
	envOld := "COPY --from=builder /opt/spack-environment /opt/spack-environment"
	envNew := "ENV PATH=/opt/view/bin:/opt/spack/bin:$PATH\n\n" + envOld
	dockerfile = strings.Replace(dockerfile, envOld, envNew, 1)

	// Modify entrypoint
	entrypointOld := `ENTRYPOINT ["/bin/bash", "--rcfile", "/etc/profile", "-l"]`
	entrypointNew := `ENTRYPOINT ["/bin/bash", "--rcfile", "/etc/profile", "-l", "-c"]`
	dockerfile = strings.Replace(dockerfile, entrypointOld, entrypointNew, 1)

	// Add Autamus Repo to Container
	addHook := "as builder"
	addCommand := "as builder\n\nADD repo /opt/spack/var/spack/repos/builtin/packages/"
	dockerfile = strings.Replace(dockerfile, addHook, addCommand, 1)

	// Add support for build cache
	buildOld := "RUN cd /opt/spack-environment && spack env activate . && spack install --fail-fast && spack gc -y"
	buildPublish := "RUN --mount=type=secret,id=aws_id " +
		"--mount=type=secret,id=aws_secret " +
		"--mount=type=secret,id=sign_key " +
		"--mount=type=secret,id=package " +
		"cat /run/secrets/package " +
		"&& cd /opt/spack-environment && spack env activate . " +
		"&& export AWS_ACCESS_KEY_ID=$(cat /run/secrets/aws_id) " +
		"&& export AWS_SECRET_ACCESS_KEY=$(cat /run/secrets/aws_secret) " +
		"&& curl http://s3.amazonaws.com/autamus-cache/build_cache/_pgp/FFEB24B0A9D81F6D5597F9900B59588C86C41BE7.pub > key.pub " +
		"&& spack gpg trust key.pub && spack install --fail-fast " +
		"&& spack gpg trust /run/secrets/sign_key " +
		"&& spack buildcache create -m autamus && spack gc -y"
	buildPR := "RUN cd /opt/spack-environment && spack env activate . " +
		"&& curl http://s3.amazonaws.com/autamus-cache/build_cache/_pgp/FFEB24B0A9D81F6D5597F9900B59588C86C41BE7.pub > key.pub " +
		"&& spack gpg trust key.pub && spack install --fail-fast && spack gc -y"

	if isPR {
		dockerfile = strings.Replace(dockerfile, buildOld, buildPR, 1)
	} else {
		dockerfile = strings.Replace(dockerfile, buildOld, buildPublish, 1)
	}

	return dockerfile, nil
}
