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
func Containerize(sEnv repo.SpackEnv, isPR bool, PublicKeyURL string) (dockerfile string, err error) {
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
	buildOldMonitor := "RUN  cd /opt/spack-environment && spack env activate . && export SPACKMON_USER=$(cat /run/secrets/su) && export SPACKMON_TOKEN=$(cat /run/secrets/st) && spack install --fail-fast && spack gc -y"
	buildPublish := "RUN --mount=type=secret,id=aws_id " +
		"--mount=type=secret,id=aws_secret " +
		"--mount=type=secret,id=sign_key " +
		"cd /opt/spack-environment && spack env activate . " +
		"&& export AWS_ACCESS_KEY_ID=$(cat /run/secrets/aws_id) " +
		"&& export AWS_SECRET_ACCESS_KEY=$(cat /run/secrets/aws_secret) " +
		"&& curl " + PublicKeyURL + " > key.pub " +
		"&& spack gpg trust key.pub " +
		"&& spack install --fail-fast" +
		"&& spack gpg trust /run/secrets/sign_key " +
		"&& spack buildcache create -r -a -m autamus && spack gc -y"

	buildPR := "RUN cd /opt/spack-environment " +
		"&& spack env activate . " +
		"&& curl " + PublicKeyURL + " > key.pub " +
		"&& spack gpg trust key.pub " +
		"&& spack install --fail-fast --monitor --monitor-save-local --monitor-tags autamus; " +
		`stat=$?; cd ~/.spack/reports/monitor/; tar -czvf monitor.tar.gz *; ` +
		`curl -F "upload=@monitor.tar.gz" http://localhost:4500/upload; exit $stat ` +
		"&& spack gc -y; "

	if len(sEnv.Spack.Mirrors) > 0 {
		if isPR {
			dockerfile = strings.Replace(dockerfile, buildOld, buildPR, 1)
			dockerfile = strings.Replace(dockerfile, buildOldMonitor, buildPR, 1)
		} else {
			dockerfile = strings.Replace(dockerfile, buildOld, buildPublish, 1)
			dockerfile = strings.Replace(dockerfile, buildOldMonitor, buildPublish, 1)
		}
	}

	return dockerfile, nil
}
