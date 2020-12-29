package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/juju/errors"
	"github.com/rs/zerolog/log"
)

// Run executes anycommand in workDir and returns stdout and error.
func Run(path string, workDir string, args ...string) (string, error) {
	_, err := exec.LookPath(path)
	if err != nil {
		log.Error().Msgf("%s required, please install it before using this tool", path)
		os.Exit(1)
	}
	cmd := exec.Command(path, args...)
	cmd.Dir = workDir
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	log.Debug().Msgf("Executing " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	output := strings.TrimSpace(stdout.String())
	errOut := strings.TrimSpace(stderr.String())
	if output != "" {
		log.Debug().Msgf(output)
	}
	if errOut != "" {
		log.Debug().Msgf(errOut)
	}
	if _, ok := err.(*exec.ExitError); !ok {
		return output, errors.Trace(err)
	}
	if err != nil {
		return output, errors.Annotatef(err, path, cmd.Args, errOut)
	}
	return output, nil
}
